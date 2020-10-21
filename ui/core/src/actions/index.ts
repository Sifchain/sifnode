// import addLiquidity from "./addLiquidity";
// import broadcastTx from "./broadcastTx";
// import connectToWallet from "./connectToWallet";
// import createPool from "./createPool";
// import destroyPool from "./destroyPool";
// import removeLiquidity from "./removeLiquidity";
// import queryListOfAvailableTokens from "./queryListOfAvailableTokens";
// import setQuantityOfToken from "./setQuantityOfToken";
// import swapTokens from "./swapTokens";
import { Api, WithApi } from "../api";
import { Store, WithStore } from "../store";
import ethWalletActions from "./ethWalletActions";
import sifWalletActions from "./sifWalletActions";
import tokenActions from "./tokenActions";

export type ActionContext<
  T extends keyof Api = keyof Api,
  S extends keyof Store = keyof Store
> = WithApi<T> & WithStore<S>;

export function createActions(context: ActionContext) {
  return {
    ...ethWalletActions(context),
    ...sifWalletActions(context),
    ...tokenActions(context),
    // ...broadcastTx(context),
    // ...connectToWallet(context),
    // ...createPool(context),
    // ...destroyPool(context),
    // ...removeLiquidity(context),
    // ...queryListOfAvailableTokens(context),
    // ...setQuantityOfToken(context),
    // ...swapTokens(context),
  };
}

export type Actions = ReturnType<typeof createActions>;
