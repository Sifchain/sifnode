import { IAmount } from "./Amount";
export declare function slipAdjustment(r: IAmount, // Native amount added
a: IAmount, // External amount added
R: IAmount, // Native Balance (before)
A: IAmount, // External Balance (before)
P: IAmount): IAmount;
/**
 *
 * @param r Native amount added
 * @param a External amount added
 * @param R Native Balance (before)
 * @param A External Balance (before)
 * @param P Existing Pool Units
 * @returns
 */
export declare function calculatePoolUnits(r: IAmount, // Native amount added
a: IAmount, // External amount added
R: IAmount, // Native Balance (before)
A: IAmount, // External Balance (before)
P: IAmount): IAmount;
export declare function calculateWithdrawal({ poolUnits, nativeAssetBalance, externalAssetBalance, lpUnits, wBasisPoints, asymmetry, }: {
    poolUnits: IAmount;
    nativeAssetBalance: IAmount;
    externalAssetBalance: IAmount;
    lpUnits: IAmount;
    wBasisPoints: IAmount;
    asymmetry: IAmount;
}): {
    withdrawNativeAssetAmount: IAmount;
    withdrawExternalAssetAmount: IAmount;
    lpUnitsLeft: IAmount;
    swapAmount: IAmount;
};
/**
 * Calculate Swap Result based on formula ( x * X * Y ) / ( x + X ) ^ 2
 * @param X  External Balance
 * @param x Swap Amount
 * @param Y Native Balance
 * @returns swapAmount
 */
export declare function calculateSwapResult(x: IAmount, X: IAmount, Y: IAmount): IAmount;
export declare function calculateExternalExternalSwapResult(ax: IAmount, // Swap Amount
aX: IAmount, // External Balance
aY: IAmount, // Native Balance
bX: IAmount, // External Balance
bY: IAmount): IAmount;
export declare function calculateReverseSwapResult(S: IAmount, X: IAmount, Y: IAmount): IAmount;
/**
 * Calculate Provider Fee according to the formula: ( x^2 * Y ) / ( x + X )^2
 * @param x Swap Amount
 * @param X External Balance
 * @param Y Native Balance
 * @returns providerFee
 */
export declare function calculateProviderFee(x: IAmount, X: IAmount, Y: IAmount): IAmount;
/**
 * Calculate price impact according to the formula (x) / (x + X)
 * @param x Swap Amount
 * @param X External Balance
 * @returns
 */
export declare function calculatePriceImpact(x: IAmount, X: IAmount): IAmount;
