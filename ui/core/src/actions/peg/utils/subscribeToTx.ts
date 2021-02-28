import { reactive } from "@vue/reactivity";
import { PegTxEventEmitter } from "../../../api/EthbridgeService/PegTxEventEmitter";
import notify from "../../../api/utils/Notifications";
import { TransactionStatus } from "../../../entities";
import { WithStore } from "../../../store";

// Using PascalCase to signify this is a factory
export function SubscribeToTx({ store }: WithStore<"tx" | "wallet">) {
  // Helper to set store tx status
  // Should this live behind a store service API?
  function storeSetTxStatus(
    hash: string | undefined,
    state: TransactionStatus
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

      notify({
        type: "info",
        message: "Pegged Transaction Pending",
        detail: {
          type: "etherscan",
          message: txHash,
        },
        loader: true,
      });
    });

    tx.onComplete(({ txHash }) => {
      storeSetTxStatus(txHash, {
        hash: txHash,
        memo: "Transaction Complete",
        state: "completed",
      });

      notify({
        type: "success",
        message: `Transfer ${txHash} has succeded.`,
      });

      // tx is complete so we can unsubscribe
      unsubscribe();
    });

    tx.onError(err => {
      storeSetTxStatus(tx.hash, {
        hash: tx.hash!, // wont matter if tx.hash doesnt exist
        memo: "Transaction Failed",
        state: "failed",
      });
      notify({ type: "error", message: err.payload.memo! });
    });

    // HACK: Trigger all hashHandlers
    // Maybe make this some kind of ready function?
    if (tx.hash) tx.setTxHash(tx.hash);

    return unsubscribe;
  };
}
