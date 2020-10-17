import { ActionContext } from "..";

export default ({ api, store }: ActionContext<"TokenService", "asset">) => ({
  async refreshTokens() {
    const top20Tokens = await api.TokenService.getTop20Tokens();
    store.asset.top20Tokens = top20Tokens;
  },
});
