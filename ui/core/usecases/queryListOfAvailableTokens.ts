import { Context } from ".";
import { amountToToken } from "../entities";

export default ({ api, store }: Context<"walletService">) => ({
  async updateListOfAvailableTokens() {
    const walletBalances = await api.walletService.getAssetBalances();

    store.setUserBalances(walletBalances);
  },
});
