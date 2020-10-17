import { ActionContext } from "..";
import { USDT, BNB } from "../constants/tokens";
export default ({ api, store }: ActionContext<"TokenService", "asset">) => ({
  async refreshTokens() {
    const top20Tokens = await api.TokenService.getTop20Tokens();
    store.asset.top20Tokens = [USDT, BNB, ...top20Tokens];
  },
});
