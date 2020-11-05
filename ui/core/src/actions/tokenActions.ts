import { ActionContext } from "..";

export default ({ api, store }: ActionContext<"TokenService", "asset">) => {
  const actions = {
    async refreshTokens() {
      // store.asset.topTokens = await api.TokenService.getTopAssets();
    },
  };

  actions.refreshTokens();

  return actions;
};
