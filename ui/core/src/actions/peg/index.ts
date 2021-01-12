import { ActionContext } from "..";
import { AssetAmount } from "../../entities";
import notify from "../../api/utils/Notifications";

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
    async burn(assetAmount: AssetAmount) {
      const tx = await api.EthbridgeService.burn({
        assetAmount,
        ethereumRecipient: store.wallet.eth.address,
        fromAddress: store.wallet.sif.address,
      });

      return await api.SifService.signAndBroadcast(tx.value.msg);
    },
    async lock(assetAmount: AssetAmount) {
      return await new Promise<any>((done) => {
        // listen for 50 confirmations
        // Eventually this should be set on ebrelayer
        // to centralize the business logic
        api.EthbridgeService.lock(store.wallet.sif.address, assetAmount, 50)
          .onTxHash(done)
          .onError((err) => {
            const payload: any = err.payload;
            notify({ type: "error", message: payload.message ?? err });
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
