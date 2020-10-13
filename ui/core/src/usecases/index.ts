import addLiquidity from "./addLiquidity";
import broadcastTx from "./broadcastTx";
import connectToWallet from "./connectToWallet";
import createPool from "./createPool";
import destroyPool from "./destroyPool";
import removeLiquidity from "./removeLiquidity";
import queryListOfAvailableTokens from "./queryListOfAvailableTokens";
import setQuantityOfToken from "./setQuantityOfToken";
import swapTokens from "./swapTokens";

import { FullApi, Api } from "../api";
import { State, Store } from "../store";

export type Context<T extends keyof FullApi = keyof FullApi> = Api<
  T,
  { state: State; store: Store }
>;

export function createUsecases(context: Context) {
  return {
    ...addLiquidity(context),
    ...broadcastTx(context),
    ...connectToWallet(context),
    ...createPool(context),
    ...destroyPool(context),
    ...removeLiquidity(context),
    ...queryListOfAvailableTokens(context),
    ...setQuantityOfToken(context),
    ...swapTokens(context),
  };
}

export type UseCases = ReturnType<typeof createUsecases>;
