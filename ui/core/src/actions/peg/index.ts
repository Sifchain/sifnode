import { ActionContext } from "..";
import {
  Address,
  Asset,
  AssetAmount,
  TransactionStatus,
} from "../../entities";
import notify from "../../api/utils/Notifications";
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
  "SifService" | "EthbridgeService" | "EthereumService",
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
        ? "18332015000000000"
        : "16164980000000000";

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

      console.log("unpeg", tx, assetAmount, store.wallet.eth.address, store.wallet.sif.address, feeAmount);

      const txStatus = await api.SifService.signAndBroadcast(tx.value.msg);

      if (txStatus.state !== "accepted") {
        notify({
          type: "error",
          message: txStatus.memo || "There was an error while unpegging",
        });
      }
      console.log("unpeg txStatus.state", txStatus.state, txStatus.memo, txStatus.code, tx.value.msg);

      return txStatus;
    },
    async approve(address: Address ,assetAmount: AssetAmount) {
       return await api.EthbridgeService.approveBridgeBankSpend(address, assetAmount)
    },
    async peg(assetAmount: AssetAmount) {
      const lockOrBurnFn = isOriginallySifchainNativeToken(assetAmount.asset)
        ? api.EthbridgeService.burnToSifchain
        : api.EthbridgeService.lockToSifchain;

      return await new Promise<TransactionStatus>(done => {
        lockOrBurnFn(store.wallet.sif.address, assetAmount, ETH_CONFIRMATIONS)
          .onTxHash(hash =>
            done({
              hash: hash.txHash,
              memo: "Transaction Accepted",
              state: "accepted",
            })
          )
          .onError(err => {
            notify({ type: "error", message: err.payload.memo! });
            done(err.payload);
          })
          .onComplete(({ txHash }) => {
            notify({
              type: "success",
              message: `Transfer ${txHash} has succeded.`,
            });
          });
      });
    },
  };

  return actions;
};
