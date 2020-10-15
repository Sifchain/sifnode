// import { Asset, Token, TokenAmount } from "../../entities";
import { ActionContext } from "..";
// import { BigintIsh } from "src/entities/fraction/Fraction";
// import JSBI from "jsbi";

// function toTokenAmount(amount: BigintIsh) {
//   return (token: Token) => {
//     return TokenAmount.create(token, amount);
//   };
// }

// const notInAssetList = (assets: Asset[]) => (asset: Asset) => {
//   return !assets.find(({ symbol }) => symbol === asset.symbol);
// };

export default ({
  api,
  store,
}: ActionContext<"walletService", "wallet" | "asset">) => ({
  // async fetchAvailableTokens() {
  //   const walletBalances = await api.walletService.getAssetBalances({
  //     limit: 10,
  //   });
  //   const topERCTokens = await api.tokenService.getTopERC20Tokens({
  //     limit: 20,
  //   });
  //   const walletTokens = walletBalances.map((assetAmount) => assetAmount.asset);
  //   const availableEmptyTokenAccounts = topERCTokens
  //     .filter(notInAssetList(walletTokens))
  //     .map(toTokenAmount(JSBI.BigInt("0")));
  //   store.setTokenBalances([...walletBalances, ...availableEmptyTokenAccounts]);
  // },
  // async refreshTokenList() {
  //   // const tokens = await api.tokenService.getSupportedTokens();
  //   // for (let token of tokens) {
  //   //   store.asset.assetMap.set(token.symbol, token);
  //   // }
  // },
  async disconnectWallet() {
    await api.walletService.disconnect();
    store.wallet.isConnected = false;
    store.wallet.balances = [];
  },
  async connectToWallet() {
    await api.walletService.connect();
    store.wallet.isConnected = api.walletService.isConnected();
    await this.refreshWalletBalances();
    // How do we listen for connection status?
    // What happens when the wallet drops off?
  },
  async refreshWalletBalances() {
    const balances = await api.walletService.getBalance();
    store.wallet.balances = balances;
  },
});
