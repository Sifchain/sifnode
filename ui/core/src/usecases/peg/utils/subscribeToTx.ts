import { reactive } from "@vue/reactivity";
import { UsecaseContext } from "../..";
import { PegTxEventEmitter } from "../../../services/EthbridgeService/PegTxEventEmitter";
import { TransactionStatus } from "../../../entities";

// Using PascalCase to signify this is a factory
export function SubscribeToTx({
  services,
  store,
}: UsecaseContext<"bus", "wallet" | "tx">) {
  // Helper to set store tx status
  // Should this live behind a store service API?
  function storeSetTxStatus(
    hash: string | undefined,
    state: TransactionStatus,
  ) {
    if (!hash || !store.wallet.eth.address) return;

    store.tx.eth[store.wallet.eth.address] =
      store.tx.eth[store.wallet.eth.address] || reactive({});

    store.tx.eth[store.wallet.eth.address][hash] = state;
  }

  /**
   * Track changes to a tx emitter send notifications
   * and update a key in the store
   * @param tx with hash set
   */
  return function subscribeToTx(tx: PegTxEventEmitter) {
    function unsubscribe() {
      tx.removeListeners();
    }

    tx.onTxHash(({ txHash }) => {
      storeSetTxStatus(txHash, {
        hash: txHash,
        memo: "Transaction Accepted",
        state: "accepted",
        symbol: tx.symbol,
      });

      services.bus.dispatch({
        type: "PegTransactionPendingEvent",
        payload: {
          hash: txHash,
        },
      });
    });

    tx.onComplete(({ txHash }) => {
      storeSetTxStatus(txHash, {
        hash: txHash,
        memo: "Transaction Complete",
        state: "completed",
      });

      services.bus.dispatch({
        type: "PegTransactionCompletedEvent",
        payload: {
          hash: txHash,
        },
      });

      // tx is complete so we can unsubscribe
      unsubscribe();
    });

    tx.onError((err) => {
      storeSetTxStatus(tx.hash, {
        hash: tx.hash || "",
        memo: "Transaction Failed",
        state: "failed",
      });

      services.bus.dispatch({
        type: "PegTransactionErrorEvent",
        payload: {
          txStatus: {
            hash: tx.hash || "",
            memo: "Transaction Error",
            state: "failed",
          },
          message: err.payload.memo!,
        },
      });
    });

    // HACK: Trigger all hashHandlers
    // Maybe make this some kind of ready function?
    if (tx.hash) tx.setTxHash(tx.hash);

    return unsubscribe;
  };
}
