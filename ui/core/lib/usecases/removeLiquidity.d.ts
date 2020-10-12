import { Context } from ".";
import { Pair, Token, TokenAmount } from "../entities";
declare function renderRemoveLiquidityPageData(liquidityPool: Pair, token: Token, tokenAmount: TokenAmount): {
    canRemoveLiquidity: boolean;
    amount: TokenAmount;
    gasFees: TokenAmount;
    shareOfPool: number;
    amountToRemoveIsTooHigh: boolean;
};
declare const _default: ({ api, store }: Context) => {
    intializeRemoveLiquidity(): void;
    renderRemoveLiquidityPageData: typeof renderRemoveLiquidityPageData;
    removeLiquidity(liquidityPool: Pair, token: Token, tokenAmount: TokenAmount): Promise<void>;
};
export default _default;
