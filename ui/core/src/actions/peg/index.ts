import { ActionContext } from "..";
import { Address, Asset, AssetAmount, TransactionStatus } from "../../entities";
import JSBI from "jsbi";

function isOriginallySifchainNativeToken(asset: Asset) {
  return ["erowan", "rowan"].includes(asset.symbol);
}
// listen for 50 confirmations
// Eventually this should be set on ebrelayer
// to centralize the business logic
const ETH_CONFIRMATIONS = 50;

export default ({
  api,
  store,
}: ActionContext<
  "SifService" | "EthbridgeService" | "NotificationService" | "EthereumService",
  "wallet"
>) => {
  const actions = {
    getSifTokens() {
      return api.SifService.getSupportedTokens();
    },

    getEthTokens() {
      return api.EthereumService.getSupportedTokens();
    },

    calculateUnpegFee(asset: Asset) {
      const feeNumber = isOriginallySifchainNativeToken(asset)
        ? "54080000000000000"
        : "58560000000000000";

      return AssetAmount(Asset.get("ceth"), JSBI.BigInt(feeNumber), {
        inBaseUnit: true,
      });
    },

    async unpeg(assetAmount: AssetAmount) {
      const lockOrBurnFn = isOriginallySifchainNativeToken(assetAmount.asset)
        ? api.EthbridgeService.lockToEthereum
        : api.EthbridgeService.burnToEthereum;

      const feeAmount = this.calculateUnpegFee(assetAmount.asset);

      const tx = await lockOrBurnFn({
        assetAmount,
        ethereumRecipient: store.wallet.eth.address,
        fromAddress: store.wallet.sif.address,
        feeAmount,
      });

      console.log(
        "unpeg",
        tx,
        assetAmount,
        store.wallet.eth.address,
        store.wallet.sif.address,
        feeAmount
      );

      const txStatus = await api.SifService.signAndBroadcast(tx.value.msg);

      if (txStatus.state !== "accepted") {
        api.NotificationService.notify({
          type: "PegTransactionErrorEvent",
          payload: {
            txStatus,
            message: txStatus.memo || "There was an error while unpegging",
          },
        });
      }

      console.log(
        "unpeg txStatus.state",
        txStatus.state,
        txStatus.memo,
        txStatus.code,
        tx.value.msg
      );

      return txStatus;
    },
    // TODO: Move this approval command to within peg and report status via store or some other means
    //       This has been done for convenience but we should not have to know in the view that
    //       approval is required before pegging as that is very much business domain knowledge
    async approve(address: Address, assetAmount: AssetAmount) {
      return await api.EthbridgeService.approveBridgeBankSpend(
        address,
        assetAmount
      );
    },
    async peg(assetAmount: AssetAmount) {
      const lockOrBurnFn = isOriginallySifchainNativeToken(assetAmount.asset)
        ? api.EthbridgeService.burnToSifchain
        : api.EthbridgeService.lockToSifchain;
      return await new Promise<TransactionStatus>(done => {
        let txHash: string = "";
        lockOrBurnFn(store.wallet.sif.address, assetAmount, ETH_CONFIRMATIONS)
          .onTxHash(hash => {
            txHash = hash.txHash;
            // TODO: Set tx status on store for pending txs to use elsewhere
            api.NotificationService.notify({
              type: "PegTransactionPendingEvent",
              payload: {
                hash: txHash,
              },
              // message: "Pegged Transaction Pending",
              // detail: {
              //   type: "etherscan",
              //   message: hash.txHash,
              // },
              // loader: true,
            });

            done({
              hash: txHash,
              memo: "Transaction Accepted",
              state: "accepted",
            });
          })
          .onError(err => {
            api.NotificationService.notify({
              type: "PegTransactionErrorEvent",
              payload: {
                txStatus: {
                  hash: txHash,
                  memo: "Transaction Accepted",
                  state: "accepted",
                },
                message: err.payload.memo!,
              },
            });
            done(err.payload);
          })
          .onComplete(({ txHash }) => {
            api.NotificationService.notify({
              type: "PegTransactionCompletedEvent",
              payload: {
                hash: txHash,
              },
              // message: `Transfer ${txHash} has succeded.`,
            });
          });
      });
    },
  };

  return actions;
};
