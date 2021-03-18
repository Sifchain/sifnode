import Big from "big.js";
import { Fraction, IFraction } from "./fraction/Fraction";

export function slipAdjustment(
  r: IFraction, // Native amount added
  a: IFraction, // External amount added
  R: IFraction, // Native Balance (before)
  A: IFraction, // External Balance (before)
  P: IFraction, // existing Pool Units
) {
  // slipAdjustment = ((R a - r A)/((r + R) (a + A)))
  const slipAdjDenominator = r.add(R).multiply(a.add(A));
  let slipAdjustmentReciprocal: IFraction;
  if (R.multiply(a).greaterThan(r.multiply(A))) {
    slipAdjustmentReciprocal = R.multiply(a)
      .subtract(r.multiply(A))
      .divide(slipAdjDenominator);
  } else {
    slipAdjustmentReciprocal = r
      .multiply(A)
      .subtract(R.multiply(a))
      .divide(slipAdjDenominator);
  }
  console.log("slipAdjustmentReciprocal", slipAdjustmentReciprocal.toFixed(4));
  // (1 - ABS((R a - r A)/((2 r + R) (a + A))))
  return new Fraction("1").subtract(slipAdjustmentReciprocal);
}

/**
 *
 * @param r Native amount added
 * @param a External amount added
 * @param R Native Balance (before)
 * @param A External Balance (before)
 * @param P Existing Pool Units
 * @returns
 */
export function calculatePoolUnits(
  r: IFraction, // Native amount added
  a: IFraction, // External amount added
  R: IFraction, // Native Balance (before)
  A: IFraction, // External Balance (before)
  P: IFraction, // existing Pool Units
) {
  if (A.equalTo("0") || R.equalTo("0") || P.equalTo("0")) {
    return r;
  }

  if (a.equalTo("0") && r.equalTo("0")) {
    return new Fraction("0");
  }

  const slipAdjustmentCalc = slipAdjustment(r, a, R, A, P);

  // ((P (a R + A r))
  const numerator = P.multiply(a.multiply(R).add(A.multiply(r)));
  const denominator = new Fraction("2").multiply(A).multiply(R);

  return numerator.divide(denominator).multiply(slipAdjustmentCalc);
}

function abs(num: Fraction) {
  if (num.lessThan("0")) {
    return num.multiply("-1");
  }
  return num;
}

const TEN_THOUSAND = new Fraction("10000");

export function calculateWithdrawal({
  poolUnits,
  nativeAssetBalance,
  externalAssetBalance,
  lpUnits,
  wBasisPoints,
  asymmetry,
}: {
  poolUnits: IFraction;
  nativeAssetBalance: IFraction;
  externalAssetBalance: IFraction;
  lpUnits: IFraction;
  wBasisPoints: IFraction;
  asymmetry: IFraction;
}) {
  let unitsToClaim = new Fraction("0");
  if (!wBasisPoints.equalTo("0")) {
    unitsToClaim = lpUnits.divide(TEN_THOUSAND.divide(wBasisPoints));
  }

  let poolUnitsOverUnitsToClaim = new Fraction("0");
  if (!unitsToClaim.equalTo("0")) {
    poolUnitsOverUnitsToClaim = poolUnits.divide(unitsToClaim);
  }

  let withdrawExternalAssetAmountPreSwap = new Fraction("0");
  let withdrawNativeAssetAmountPreSwap = new Fraction("0");
  if (!poolUnitsOverUnitsToClaim.equalTo("0")) {
    withdrawExternalAssetAmountPreSwap = externalAssetBalance.divide(
      poolUnitsOverUnitsToClaim,
    );

    withdrawNativeAssetAmountPreSwap = nativeAssetBalance.divide(
      poolUnitsOverUnitsToClaim,
    );
  }

  const lpUnitsLeft = lpUnits.subtract(unitsToClaim);

  const swapAmount = abs(
    asymmetry.equalTo("0")
      ? new Fraction("0")
      : asymmetry.lessThan("0")
      ? externalAssetBalance.divide(
          poolUnits.divide(unitsToClaim.divide(TEN_THOUSAND.divide(asymmetry))),
        )
      : nativeAssetBalance.divide(
          poolUnits.divide(unitsToClaim.divide(TEN_THOUSAND.divide(asymmetry))),
        ),
  );

  const newExternalAssetBalance = externalAssetBalance.subtract(
    withdrawExternalAssetAmountPreSwap,
  );

  const newNativeAssetBalance = nativeAssetBalance.subtract(
    withdrawNativeAssetAmountPreSwap,
  );

  const withdrawNativeAssetAmount = !asymmetry.lessThan("0")
    ? withdrawNativeAssetAmountPreSwap.subtract(swapAmount)
    : withdrawNativeAssetAmountPreSwap.add(
        calculateSwapResult(
          abs(swapAmount),
          newExternalAssetBalance,
          newNativeAssetBalance,
        ),
      );

  const withdrawExternalAssetAmount = asymmetry.lessThan("0")
    ? withdrawExternalAssetAmountPreSwap.subtract(swapAmount)
    : withdrawExternalAssetAmountPreSwap.add(
        calculateSwapResult(
          abs(swapAmount),
          newNativeAssetBalance,
          newExternalAssetBalance,
        ),
      );

  return {
    withdrawNativeAssetAmount,
    withdrawExternalAssetAmount,
    lpUnitsLeft,
    swapAmount,
  };
}

/**
 * Calculate Swap Result based on formula ( x * X * Y ) / ( x + X ) ^ 2
 * @param X  External Balance
 * @param x Swap Amount
 * @param Y Native Balance
 * @returns swapAmount
 */
export function calculateSwapResult(x: IFraction, X: IFraction, Y: IFraction) {
  if (x.equalTo("0") || X.equalTo("0") || Y.equalTo("0")) {
    return new Fraction("0");
  }
  const xPlusX = x.add(X);
  return x.multiply(X).multiply(Y).divide(xPlusX.multiply(xPlusX));
}

export function calculateExternalExternalSwapResult(
  // External -> Native pool
  ax: IFraction, // Swap Amount
  aX: IFraction, // External Balance
  aY: IFraction, // Native Balance
  // Native -> External pool
  bX: IFraction, // External Balance
  bY: IFraction, // Native Balance
) {
  const ay = calculateSwapResult(ax, aX, aY);
  return calculateSwapResult(ay, bX, bY);
}

// Formula: S = (x * X * Y) / (x + X) ^ 2
// Reverse Formula: x = ( -2*X*S + X*Y - X*sqrt( Y*(Y - 4*S) ) ) / 2*S
// Need to use Big.js for sqrt calculation
// Ok to accept a little precision loss as reverse swap amount can be rough
export function calculateReverseSwapResult(S: Big, X: Big, Y: Big) {
  // Adding a check here because sqrt of a negative number will throw an exception
  if (S.eq("0") || X.eq("0") || S.times(4).gt(Y)) {
    return Big("0");
  }
  const term1 = Big(-2).times(X).times(S);
  const term2 = X.times(Y);
  const underRoot = Y.times(Y.minus(S.times(4)));
  const term3 = X.times(underRoot.sqrt());
  const numerator = term1.plus(term2).minus(term3);
  const denominator = S.times(2);
  const x = numerator.div(denominator);
  return x.gte(Big("0")) ? x : Big("0");
}

/**
 * Calculate Provider Fee according to the formula: ( x^2 * Y ) / ( x + X )^2
 * @param x Swap Amount
 * @param X External Balance
 * @param Y Native Balance
 * @returns providerFee
 */
export function calculateProviderFee(x: IFraction, X: IFraction, Y: IFraction) {
  if (x.equalTo("0") || X.equalTo("0") || Y.equalTo("0")) {
    return new Fraction("0");
  }
  const xPlusX = x.add(X);
  return x.multiply(x).multiply(Y).divide(xPlusX.multiply(xPlusX));
}

/**
 * Calculate price impact according to the formula (x) / (x + X)
 * @param x Swap Amount
 * @param X External Balance
 * @returns
 */
export function calculatePriceImpact(x: IFraction, X: IFraction) {
  if (x.equalTo("0")) {
    return new Fraction("0");
  }
  const denominator = x.add(X);
  return x.divide(denominator);
}
