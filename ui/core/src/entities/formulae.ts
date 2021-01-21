import Big from "big.js";
import {AssetAmount, IAssetAmount} from "./AssetAmount";
import {Fraction, IFraction} from "./fraction/Fraction";
import {IPool} from "./Pool";
import JSBI from "jsbi";

export function calcLpUnits(
  amounts: [IAssetAmount, IAssetAmount],
  nativeAssetAmount: AssetAmount,
  externalAssetAmount: AssetAmount
) {
  // Not necessarily native but we will treat it like so as the formulae are symmetrical
  const nativeAssetBalance = amounts.find(
    a => a.asset.symbol === nativeAssetAmount.asset.symbol
  );
  const externalAssetBalance = amounts.find(
    a => a.asset.symbol === externalAssetAmount.asset.symbol
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
  if (x.equalTo("0") || Y.equalTo("0")) return new Fraction("0");
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

// Formula: ( x^2 * Y ) / ( x + X )^2
export function calculateProviderFee(x: IFraction, X: IFraction, Y: IFraction) {
  if (x.equalTo("0") || Y.equalTo("0")) return new Fraction("0");
  const xPlusX = x.add(X);
  return x
    .multiply(x)
    .multiply(Y)
    .divide(xPlusX.multiply(xPlusX));
}

// (x) / (x + X)
export function calculatePriceImpact(x: IFraction, X: IFraction) {
  if (x.equalTo("0") || X.equalTo("0")) return new Fraction("0");
  const denominator = x.add(X);
  return x.divide(denominator);
}

export const getSwapSlip = (x: AssetAmount, pool: IPool, toRowan: boolean): AssetAmount => {
  // formula: (x) / (x + X)
  const X = pool.amounts.find(a => (toRowan ? a.asset.symbol === x.asset.symbol : a.asset.symbol !== x.asset.symbol));
  if (x && X && JSBI.toNumber(x.amount) && JSBI.toNumber(X.amount)) {
    return AssetAmount(x.asset, x.divide(x.add(X)));
  } else {
    return AssetAmount(x.asset, "0");
  }
};

export const getSwapFee = (x: AssetAmount, pool: IPool, toRowan: boolean): AssetAmount => {
  // formula: (x * x * Y) / (x + X) ^ 2
  const X = pool.amounts.find(a => (toRowan ? a.asset.symbol === x.asset.symbol : a.asset.symbol !== x.asset.symbol));
  const Y = pool.amounts.find(a => (toRowan ? a.asset.symbol !== x.asset.symbol : a.asset.symbol === x.asset.symbol));
  if (x && X && Y && JSBI.toNumber(x.amount) && JSBI.toNumber(X.amount)) {
    const numerator = x.multiply(x).multiply(Y);
    const xPlusX = x.add(X);
    const denominator = xPlusX.multiply(xPlusX);
    return AssetAmount(Y.asset, numerator.divide(denominator));
  } else {
    return AssetAmount(x.asset, "0");
  }
};

export const getSwapOutput = (x: AssetAmount, pool: IPool, toRowan: boolean): AssetAmount => {
  // formula: (x * X * Y) / (x + X) ^ 2
  const X = pool.amounts.find(a => (toRowan ? a.asset.symbol === x.asset.symbol : a.asset.symbol !== x.asset.symbol));
  const Y = pool.amounts.find(a => (toRowan ? a.asset.symbol !== x.asset.symbol : a.asset.symbol === x.asset.symbol));
  if (x && X && Y && JSBI.toNumber(x.amount) && JSBI.toNumber(X.amount)) {
    const numerator = x.multiply(x).multiply(Y);
    const xPlusX = x.add(X);
    const denominator = xPlusX.multiply(xPlusX);
    return AssetAmount(Y.asset, numerator.divide(denominator));
  } else {
    return AssetAmount(x.asset, "0");
  }
};

export const getSwapOutputWithFee = (
  x: AssetAmount,
  pool: IPool,
  toRowan: boolean,
  transactionFee: AssetAmount = AssetAmount(x.asset,"1"),
): AssetAmount => {
  // formula: getSwapOutput() - one rowan
  const r = getSwapOutput(x, pool, toRowan);
  const a = pool.amounts.find(a => a.asset.symbol === x.asset.symbol);
  const b = pool.amounts.find(a => a.asset.symbol !== x.asset.symbol);
  if (a && b) {
    const poolAfterTransaction: IPool = toRowan ? {
      ...pool,
      amounts: [a, b]
    } : {
      ...pool,
      amounts: [a, b]
    };
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    const rowanFee = toRowan ? transactionFee : getValueOfRowanInAsset(transactionFee, poolAfterTransaction);
    return AssetAmount(r.asset, r.subtract(rowanFee));
  } else {
    return AssetAmount(x.asset, "0");
  }
};

export const getDoubleSwapSlip = (x: AssetAmount, pool1: IPool, pool2: IPool): AssetAmount => {
  // formula: getSwapSlip1(input1) + getSwapSlip2(getSwapOutput1 => input2)
  const swapSlip1 = getSwapSlip(x, pool1, true);
  const r = getSwapOutput(x, pool1, true);
  const swapSlip2 = getSwapSlip(r, pool2, false);
  return AssetAmount(x.asset, (swapSlip1.add(swapSlip2)).multiply("100"));
};

export const getDoubleSwapFee = (x: AssetAmount, pool1: IPool, pool2: IPool): AssetAmount => {
  // formula: getSwapFee1 + getSwapFee2
  const fee1 = getSwapFee(x, pool1, true);
  const r = getSwapOutput(x, pool1, true);
  const fee2 = getSwapFee(r, pool2, false);
  const assetValue = getValueOfRowanInAsset(fee1, pool2);
  return AssetAmount(x.asset, fee2.add(assetValue));
};

export const getDoubleSwapOutputWithFee = (
  x: AssetAmount,
  pool1: IPool,
  pool2: IPool,
  transactionFee: AssetAmount = AssetAmount(x.asset, 1),
): AssetAmount => {
  // formula: (getSwapOutput(pool1) => getSwapOutput(pool2)) - rowanFee
  const r = getSwapOutput(x, pool1, true);
  const output = getSwapOutput(r, pool2, false);
  const A = pool2.amounts.find(a => a.asset.symbol === x.asset.symbol);
  const R = pool2.amounts.find(a => a.asset.symbol !== x.asset.symbol);
  if (R && A) {
    const poolAfterTransaction: IPool = {
      ...pool2,
      amounts: [AssetAmount(R.asset, R.add(r)), AssetAmount(A.asset, A.subtract(output))]
    };
    const rowanFee = getValueOfRowanInAsset(transactionFee, poolAfterTransaction);
    return AssetAmount(output.asset, output.subtract(rowanFee));
  } else {
    return AssetAmount(output.asset, 0);
  }
};

export const getValueOfRowanInAsset = (r: AssetAmount, pool: IPool): AssetAmount => {
  // formula: ((r * A) / R) => A per R ($ per rowan)
  const R = pool.amounts.find(a => a.asset.symbol !== r.asset.symbol);
  const A = pool.amounts.find(a => a.asset.symbol === r.asset.symbol);
  if (A && R && JSBI.toNumber(R.amount)) {
    return AssetAmount(A.asset, r.multiply(A).divide(R));
  } else {
    return AssetAmount(r.asset,"0");
  }
};

export const getValueOfAssetInRowan = (x: AssetAmount, pool: IPool): AssetAmount => {
  // formula: ((a * R) / A) => R per A (rowan per $)
  const R = pool.amounts.find(a => a.asset.symbol !== x.asset.symbol);
  const A = pool.amounts.find(a => a.asset.symbol === x.asset.symbol);
  if (R && A && JSBI.toNumber(A.amount)) {
    return AssetAmount(R.asset, x.multiply(R).divide(A));
  } else {
    return AssetAmount(x.asset,"0");
  }
};

export const getValueOfAsset1InAsset2 = (inputAsset: AssetAmount, pool1: IPool, pool2: IPool): AssetAmount => {
  // formula: (A2 / R) * (R / A1) => A2/A1 => A2 per A1 ($ per Asset)
  const oneAsset = AssetAmount(inputAsset.asset, 1);
  const A2perR = getValueOfRowanInAsset(oneAsset, pool2);
  const RperA1 = getValueOfAssetInRowan(inputAsset, pool1);
  return AssetAmount(inputAsset.asset, A2perR.multiply(RperA1));
};

export const assetToBase = (asset: AssetAmount): AssetAmount => {
  return AssetAmount(asset.asset, asset.multiply(JSBI.BigInt(10 ** asset.asset.decimals)).toFixed(0));
};


// export const getDoubleSwapInput = (pool1: IPool, pool2: IPool, outputAmount: AssetAmount): AssetAmount => {
//   // formula: getSwapInput(pool2) => getSwapInput(pool1)
//   const y = getSwapInput(false, pool2, outputAmount);
//   return getSwapInput(true, pool1, y)
// };

// export const getSwapInput = (toRowan: boolean, pool: IPool, y: AssetAmount): AssetAmount => {
//   // formula: (((X*Y)/y - 2*X) - sqrt(((X*Y)/y - 2*X)^2 - 4*X^2))/2
//   // (part1 - sqrt(part1 - part2))/2
//   const X = pool.amounts.find(a => (toRowan ? a.asset.symbol === y.asset.symbol : a.asset.symbol !== y.asset.symbol));
//   const Y = pool.amounts.find(a => (toRowan ? a.asset.symbol !== y.asset.symbol : a.asset.symbol === y.asset.symbol));
//   if (X && Y) {
//     const part1 = X.multiply(Y).divide(y).subtract(X.multiply("2"));
//     const part2 = (X.multiply(X)).multiply("4");
//     const result = part1.subtract((part1.multiply(part1)).subtract(part2).divide('SQRT')).divide("2");
//     return AssetAmount(X.asset, result);
//   } else {
//     return AssetAmount(y.asset, "0");
//   }
// };
