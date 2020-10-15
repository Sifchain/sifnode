// import addLiquidity from "./addLiquidity";
// import broadcastTx from "./broadcastTx";
// import connectToWallet from "./connectToWallet";
// import createPool from "./createPool";
// import destroyPool from "./destroyPool";
// import removeLiquidity from "./removeLiquidity";
// import queryListOfAvailableTokens from "./queryListOfAvailableTokens";
// import setQuantityOfToken from "./setQuantityOfToken";
// import swapTokens from "./swapTokens";
import walletActions from "./walletActions";
import tokenActions from "./tokenActions";
import { Api, WithApi } from "../api";
import { Store, WithStore } from "../store";

export type ActionContext<
  T extends keyof Api = keyof Api,
  S extends keyof Store = keyof Store
> = WithApi<T> & WithStore<S>;

export function createActions(context: ActionContext) {
  return {
    ...walletActions(context),
    ...tokenActions(context),
    // ...addLiquidity(context),
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
