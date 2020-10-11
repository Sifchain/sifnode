import { Context } from ".";
import { TokenAmount } from "../entities";
declare function renderCreatePoolData(amountA: TokenAmount, amountB: TokenAmount): {
    tokenAPerBRatio: number;
    tokenBPerARatio: number;
    tokenAAmountOwned: TokenAmount;
    tokenBAmountOwned: TokenAmount;
    shareOfPool: number;
    canCreatePool: boolean;
    isInsufficientFunds: boolean;
};
declare const _default: ({ api, store }: Context) => {
    intializeCreatePoolUseCase(): void;
    renderCreatePoolData: typeof renderCreatePoolData;
    createPool(amountA: TokenAmount, amountB: TokenAmount): Promise<void>;
};
export default _default;
