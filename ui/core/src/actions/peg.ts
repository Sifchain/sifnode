import { ActionContext } from ".";

export default ({
  api,
}: ActionContext<"SifService" | "EthereumService", "asset">) => {
  const actions = {
    getSifTokens() {
      return api.SifService.getSupportedTokens();
    },
    getEthTokens() {
      return api.EthereumService.getSupportedTokens();
    },
  };

  return actions;
};
