"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.calculatePriceImpact = exports.calculateProviderFee = exports.calculateReverseSwapResult = exports.calculateExternalExternalSwapResult = exports.calculateSwapResult = exports.calculateWithdrawal = exports.calculatePoolUnits = exports.slipAdjustment = void 0;
const Amount_1 = require("./Amount");
function slipAdjustment(r, // Native amount added
a, // External amount added
R, // Native Balance (before)
A, // External Balance (before)
P) {
    // slipAdjustment = ((R a - r A)/((r + R) (a + A)))
    const slipAdjDenominator = r.add(R).multiply(a.add(A));
    let slipAdjustmentReciprocal;
    if (R.multiply(a).greaterThan(r.multiply(A))) {
        slipAdjustmentReciprocal = R.multiply(a)
            .subtract(r.multiply(A))
            .divide(slipAdjDenominator);
    }
    else {
        slipAdjustmentReciprocal = r
            .multiply(A)
            .subtract(R.multiply(a))
            .divide(slipAdjDenominator);
    }
    // (1 - ABS((R a - r A)/((2 r + R) (a + A))))
    return Amount_1.Amount("1").subtract(slipAdjustmentReciprocal);
}
exports.slipAdjustment = slipAdjustment;
/**
 *
 * @param r Native amount added
 * @param a External amount added
 * @param R Native Balance (before)
 * @param A External Balance (before)
 * @param P Existing Pool Units
 * @returns
 */
function calculatePoolUnits(r, // Native amount added
a, // External amount added
R, // Native Balance (before)
A, // External Balance (before)
P) {
    if (A.equalTo("0") || R.equalTo("0") || P.equalTo("0")) {
        return r;
    }
    if (a.equalTo("0") && r.equalTo("0")) {
        return Amount_1.Amount("0");
    }
    const slipAdjustmentCalc = slipAdjustment(r, a, R, A, P);
    // ((P (a R + A r))
    const numerator = P.multiply(a.multiply(R).add(A.multiply(r)));
    const denominator = Amount_1.Amount("2").multiply(A).multiply(R);
    return numerator.divide(denominator).multiply(slipAdjustmentCalc);
}
exports.calculatePoolUnits = calculatePoolUnits;
function abs(num) {
    if (num.lessThan("0")) {
        return num.multiply("-1");
    }
    return num;
}
const TEN_THOUSAND = Amount_1.Amount("10000");
function calculateWithdrawal({ poolUnits, nativeAssetBalance, externalAssetBalance, lpUnits, wBasisPoints, asymmetry, }) {
    let unitsToClaim = Amount_1.Amount("0");
    if (!wBasisPoints.equalTo("0")) {
        unitsToClaim = lpUnits.divide(TEN_THOUSAND.divide(wBasisPoints));
    }
    let poolUnitsOverUnitsToClaim = Amount_1.Amount("0");
    if (!unitsToClaim.equalTo("0")) {
        poolUnitsOverUnitsToClaim = poolUnits.divide(unitsToClaim);
    }
    let withdrawExternalAssetAmountPreSwap = Amount_1.Amount("0");
    let withdrawNativeAssetAmountPreSwap = Amount_1.Amount("0");
    if (!poolUnitsOverUnitsToClaim.equalTo("0")) {
        withdrawExternalAssetAmountPreSwap = externalAssetBalance.divide(poolUnitsOverUnitsToClaim);
        withdrawNativeAssetAmountPreSwap = nativeAssetBalance.divide(poolUnitsOverUnitsToClaim);
    }
    const lpUnitsLeft = lpUnits.subtract(unitsToClaim);
    const swapAmount = abs(asymmetry.equalTo("0")
        ? Amount_1.Amount("0")
        : asymmetry.lessThan("0")
            ? externalAssetBalance.divide(poolUnits.divide(unitsToClaim.divide(TEN_THOUSAND.divide(asymmetry))))
            : nativeAssetBalance.divide(poolUnits.divide(unitsToClaim.divide(TEN_THOUSAND.divide(asymmetry)))));
    const newExternalAssetBalance = externalAssetBalance.subtract(withdrawExternalAssetAmountPreSwap);
    const newNativeAssetBalance = nativeAssetBalance.subtract(withdrawNativeAssetAmountPreSwap);
    const withdrawNativeAssetAmount = !asymmetry.lessThan("0")
        ? withdrawNativeAssetAmountPreSwap.subtract(swapAmount)
        : withdrawNativeAssetAmountPreSwap.add(calculateSwapResult(abs(swapAmount), newExternalAssetBalance, newNativeAssetBalance));
    const withdrawExternalAssetAmount = asymmetry.lessThan("0")
        ? withdrawExternalAssetAmountPreSwap.subtract(swapAmount)
        : withdrawExternalAssetAmountPreSwap.add(calculateSwapResult(abs(swapAmount), newNativeAssetBalance, newExternalAssetBalance));
    return {
        withdrawNativeAssetAmount,
        withdrawExternalAssetAmount,
        lpUnitsLeft,
        swapAmount,
    };
}
exports.calculateWithdrawal = calculateWithdrawal;
/**
 * Calculate Swap Result based on formula ( x * X * Y ) / ( x + X ) ^ 2
 * @param X  External Balance
 * @param x Swap Amount
 * @param Y Native Balance
 * @returns swapAmount
 */
function calculateSwapResult(x, X, Y) {
    if (x.equalTo("0") || X.equalTo("0") || Y.equalTo("0")) {
        return Amount_1.Amount("0");
    }
    const xPlusX = x.add(X);
    return x.multiply(X).multiply(Y).divide(xPlusX.multiply(xPlusX));
}
exports.calculateSwapResult = calculateSwapResult;
function calculateExternalExternalSwapResult(
// External -> Native pool
ax, // Swap Amount
aX, // External Balance
aY, // Native Balance
// Native -> External pool
bX, // External Balance
bY) {
    const ay = calculateSwapResult(ax, aX, aY);
    return calculateSwapResult(ay, bX, bY);
}
exports.calculateExternalExternalSwapResult = calculateExternalExternalSwapResult;
// Formula: S = (x * X * Y) / (x + X) ^ 2
// Reverse Formula: x = ( -2*X*S + X*Y - X*sqrt( Y*(Y - 4*S) ) ) / 2*S
// Need to use Big.js for sqrt calculation
// Ok to accept a little precision loss as reverse swap amount can be rough
function calculateReverseSwapResult(S, X, Y) {
    // Adding a check here because sqrt of a negative number will throw an exception
    if (S.equalTo("0") ||
        X.equalTo("0") ||
        S.multiply(Amount_1.Amount("4")).greaterThan(Y)) {
        return Amount_1.Amount("0");
    }
    const term1 = Amount_1.Amount("-2").multiply(X).multiply(S);
    const term2 = X.multiply(Y);
    const underRoot = Y.multiply(Y.subtract(S.multiply(Amount_1.Amount("4"))));
    const term3 = X.multiply(underRoot.sqrt());
    const numerator = term1.add(term2).subtract(term3);
    const denominator = S.multiply(Amount_1.Amount("2"));
    const x = numerator.divide(denominator);
    return x.greaterThanOrEqual(Amount_1.Amount("0")) ? x : Amount_1.Amount("0");
}
exports.calculateReverseSwapResult = calculateReverseSwapResult;
/**
 * Calculate Provider Fee according to the formula: ( x^2 * Y ) / ( x + X )^2
 * @param x Swap Amount
 * @param X External Balance
 * @param Y Native Balance
 * @returns providerFee
 */
function calculateProviderFee(x, X, Y) {
    if (x.equalTo("0") || X.equalTo("0") || Y.equalTo("0")) {
        return Amount_1.Amount("0");
    }
    const xPlusX = x.add(X);
    return x.multiply(x).multiply(Y).divide(xPlusX.multiply(xPlusX));
}
exports.calculateProviderFee = calculateProviderFee;
/**
 * Calculate price impact according to the formula (x) / (x + X)
 * @param x Swap Amount
 * @param X External Balance
 * @returns
 */
function calculatePriceImpact(x, X) {
    if (x.equalTo("0")) {
        return Amount_1.Amount("0");
    }
    const denominator = x.add(X);
    return x.divide(denominator);
}
exports.calculatePriceImpact = calculatePriceImpact;
//# sourceMappingURL=formulae.js.map