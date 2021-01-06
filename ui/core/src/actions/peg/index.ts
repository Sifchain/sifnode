import { ActionContext } from "..";
import { AssetAmount } from "../../entities";
import notify from "../../api/utils/Notifications";

export default ({
  api,
}: ActionContext<
  "SifService" | "EthbridgeService" | "EthereumService",
  "asset"
>) => {
  const actions = {
    getSifTokens() {
      return api.SifService.getSupportedTokens();
    },
    getEthTokens() {
      return api.EthereumService.getSupportedTokens();
    },
    burn(ethereumRecipient: string, assetAmount: AssetAmount) {
      // Some random string for now
      // const txHash = "abcd1234";
      // Create an emitter
      // const e = createPegTxEventEmitter(txHash);
      // Direct that emitter through a mock sequence
      // return mockBurnSequence(e);
    },
    async lock(cosmosRecipient: string, assetAmount: AssetAmount) {
      // listen for 50 confirmations
      // Eventually this should be set on ebrelayer
      // to centralize the business logic
      api.EthbridgeService.lock(cosmosRecipient, assetAmount, 50)
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
    },
  };

  return actions;
};
