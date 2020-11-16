import Big from "big.js";
import { AssetAmount, IAssetAmount } from "./AssetAmount";
import { Fraction, IFraction } from "./fraction/Fraction";

export function calcLpUnits(
  amounts: [IAssetAmount, IAssetAmount],
  nativeAssetAmount: AssetAmount,
  externalAssetAmount: AssetAmount
) {
  // Not necessarily native but we will treat it like so as the formulae are symmetrical
  const nativeAssetBalance = amounts.find(
    (a) => a.asset.symbol === nativeAssetAmount.asset.symbol
  );
  const externalAssetBalance = amounts.find(
    (a) => a.asset.symbol === externalAssetAmount.asset.symbol
  );

  if (!nativeAssetBalance || !externalAssetBalance) {
    throw new Error("Pool does not contain given assets");
  }

  const R = nativeAssetBalance.add(nativeAssetAmount);
  const A = externalAssetBalance.add(externalAssetAmount);
  const r = nativeAssetAmount;
  const a = externalAssetAmount;
  const term1 = R.add(A); // R + A
  const term2 = r.multiply(A).add(R.multiply(a)); // r * A + R * a
  const numerator = term1.multiply(term2);
  const denominator = R.multiply(A).multiply("4");
  return numerator.divide(denominator);
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
  const unitsToClaim = lpUnits.divide(TEN_THOUSAND.divide(wBasisPoints));

  const poolUnitsOverUnitsToClaim = poolUnits.divide(unitsToClaim);

  const withdrawExternalAssetAmountPreSwap = externalAssetBalance.divide(
    poolUnitsOverUnitsToClaim
  );

  const withdrawNativeAssetAmountPreSwap = nativeAssetBalance.divide(
    poolUnitsOverUnitsToClaim
  );

  const lpUnitsLeft = lpUnits.subtract(unitsToClaim);

  const swapAmount = abs(
    asymmetry.equalTo("0")
      ? new Fraction("0")
      : asymmetry.lessThan("0")
      ? externalAssetBalance.divide(
          poolUnits.divide(unitsToClaim.divide(TEN_THOUSAND.divide(asymmetry)))
        )
      : nativeAssetBalance.divide(
          poolUnits.divide(unitsToClaim.divide(TEN_THOUSAND.divide(asymmetry)))
        )
  );

  const newExternalAssetBalance = externalAssetBalance.subtract(
    withdrawExternalAssetAmountPreSwap
  );

  const newNativeAssetBalance = nativeAssetBalance.subtract(
    withdrawNativeAssetAmountPreSwap
  );

  const withdrawNativeAssetAmount = !asymmetry.lessThan("0")
    ? withdrawNativeAssetAmountPreSwap.subtract(swapAmount)
    : withdrawNativeAssetAmountPreSwap.add(
        calculateSwapResult(
          newExternalAssetBalance,
          abs(swapAmount),
          newNativeAssetBalance
        )
      );

  const withdrawExternalAssetAmount = asymmetry.lessThan("0")
    ? withdrawExternalAssetAmountPreSwap.subtract(swapAmount)
    : withdrawExternalAssetAmountPreSwap.add(
        calculateSwapResult(
          newNativeAssetBalance,
          abs(swapAmount),
          newExternalAssetBalance
        )
      );

  return {
    withdrawNativeAssetAmount,
    withdrawExternalAssetAmount,
    lpUnitsLeft,
    swapAmount,
  };
}

export function calculateSwapResult(X: IFraction, x: IFraction, Y: IFraction) {
  return x
    .multiply(X)
    .multiply(Y)
    .divide(x.add(X).multiply(x.add(X)));
}

export function calculateExternalExternalSwapResult(
  // External -> Native pool
  ax: IFraction, // Swap Amount
  aX: IFraction, // External Balance
  aY: IFraction, // Native Balance
  // Native -> External pool
  bX: IFraction, // External Balance
  bY: IFraction // Native Balance
) {
  const emitAmount = calculateSwapResult(aX, ax, aY);
  return calculateSwapResult(bX, emitAmount, bY);
}

// Formula: S = (x * X * Y) / (x + X) ^ 2
// Reverse Formula: x = ( -2*X*S + X*Y - X*sqrt( Y*(Y - 4*S) ) ) / 2*S
// Need to use Big.js for sqrt calculation
// Ok to accept a little precision loss as reverse swap amount can be rough
export function calculateReverseSwapResult(S: Big, X: Big, Y: Big) {
  if (S.eq("0")) {
    return Big("0");
  }

  const term1 = Big(-2)
    .times(X)
    .times(S);

  const term2 = X.times(Y);
  const underRoot = Y.times(Y.minus(S.times(4)));

  const term3 = X.times(underRoot.sqrt());

  const numerator = term1.plus(term2).minus(term3);
  const denominator = S.times(2);

  const x = numerator.div(denominator);
  return x;
}
