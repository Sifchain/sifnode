import { Context } from ".";
import { TokenAmount } from "../entities";
declare function renderLiquidityData(amountA: TokenAmount, amountB?: TokenAmount): {
    tokenAPerBRatio: number;
    tokenBPerARatio: number;
    tokenAAmountOwned: TokenAmount;
    tokenBAmountOwned: TokenAmount;
    shareOfPool: number;
    isInsufficientFunds: boolean;
};
declare const _default: ({ api, store }: Context) => {
    intializeAddLiquidityUseCase(): void;
    renderLiquidityData: typeof renderLiquidityData;
    addLiquidity(amountA: TokenAmount, amountB: TokenAmount): Promise<void>;
};
export default _default;
