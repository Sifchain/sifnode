import { ActionContext } from "..";

export default ({ api, store }: ActionContext<"tokenService", "asset">) => ({
  async refreshTokens() {
    const top20Tokens = await api.tokenService.getTop20Tokens();
    store.asset.top20Tokens = top20Tokens;
  },
});
