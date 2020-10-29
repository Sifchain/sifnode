import { ActionContext } from "..";
import { USDT, BNB } from "../constants/tokens";

export default ({ api, store }: ActionContext<"TokenService", "asset">) => {
  const actions = {
    async refreshTokens() {
      const top20Tokens = await api.TokenService.getTop20Tokens();
      store.asset.top20Tokens = [USDT, BNB, ...top20Tokens];
    },
  };

  actions.refreshTokens();

  return actions;
};
