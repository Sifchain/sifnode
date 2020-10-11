import addLiquidity from "./addLiquidity";
import broadcastTx from "./broadcastTx";
import connectToWallet from "./connectToWallet";
import createPool from "./createPool";
import destroyPool from "./destroyPool";
import removeLiquidity from "./removeLiquidity";
import queryListOfAvailableTokens from "./queryListOfAvailableTokens";
import setQuantityOfToken from "./setQuantityOfToken";
import swapTokens from "./swapTokens";
import * as api from "../api";
import { store } from "../store";
export function createUsecases(context) {
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
export const usecases = createUsecases({ api, state: store.state, store });
//# sourceMappingURL=index.js.map