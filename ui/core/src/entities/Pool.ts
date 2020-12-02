import { Asset } from "./Asset";
import { AssetAmount, IAssetAmount } from "./AssetAmount";
import { Pair } from "./Pair";
import Big from "big.js";
import JSBI from "jsbi";
import { Fraction } from "./fraction/Fraction";
import { calcLpUnits, calculateReverseSwapResult } from "./formulae";

export type Pool = ReturnType<typeof Pool>;
export type IPool = Omit<Pool, "poolUnits" | "calculatePoolUnits">;

export function Pool(
  a: AssetAmount,
  b: AssetAmount,
  poolUnits: Fraction = new Fraction("0")
) {
  const pair = Pair(a, b);
  const amounts: [IAssetAmount, IAssetAmount] = pair.amounts;

  const instance = {
    amounts,

    otherAsset: pair.otherAsset,
    symbol: pair.symbol,
    contains: pair.contains,
    toString: pair.toString,
    poolUnits: calcLpUnits(
      [AssetAmount(a.asset, "0"), AssetAmount(b.asset, "0")],
      a,
      b
    ),

    priceAsset(asset: Asset) {
      return this.calcSwapResult(AssetAmount(asset, "1"));
    },

    // https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
    // Formula: swapAmount = (x * X * Y) / (x + X) ^ 2
    calcSwapResult(x: AssetAmount) {
      const X = amounts.find((a) => a.asset.symbol === x.asset.symbol);
      if (!X)
        throw new Error(
          `Sent amount with symbol ${
            x.asset.symbol
          } does not exist in this pair: ${this.toString()}`
        );
      const Y = amounts.find((a) => a.asset.symbol !== x.asset.symbol);
      if (!Y) throw new Error("Pool does not have an opposite asset."); // For Typescript's sake will probably never happen

      if (x.equalTo("0")) return AssetAmount(this.otherAsset(x.asset), "0");

      const swapAmount = x
        .multiply(X)
        .multiply(Y)
        .divide(x.add(X).multiply(x.add(X)));

      return AssetAmount(this.otherAsset(x.asset), swapAmount);
    },

    // Formula: S = (x * X * Y) / (x + X) ^ 2
    // Reverse Formula: x = ( -2*X*S + X*Y - X*sqrt( Y*(Y - 4*S) ) ) / 2*S
    calcReverseSwapResult(Sa: AssetAmount) {
      const Ya = amounts.find((a) => a.asset.symbol === Sa.asset.symbol);
      if (!Ya)
        throw new Error(
          `Sent amount with symbol ${
            Sa.asset.symbol
          } does not exist in this pair: ${this.toString()}`
        );
      const Xa = amounts.find((a) => a.asset.symbol !== Sa.asset.symbol);
      if (!Xa) throw new Error("Pool does not have an opposite asset."); // For Typescript's sake will probably never happen
      const otherAsset = this.otherAsset(Sa.asset);
      if (Sa.equalTo("0")) {
        return AssetAmount(otherAsset, "0");
      }

      // Need to use Big.js for sqrt calculation
      // Ok to accept a little precision loss as reverse swap amount can be rough
      const S = Big(Sa.toFixed());
      const X = Big(Xa.toFixed());
      const Y = Big(Ya.toFixed());
      const x = calculateReverseSwapResult(S, X, Y);

      return AssetAmount(otherAsset, x.toFixed());
    },
    // https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
    // poolerUnits = ((R + A) * (r * A + R * a))/(4 * R * A)
    calculatePoolUnits(
      nativeAssetAmount: AssetAmount,
      externalAssetAmount: AssetAmount
    ) {
      const lpUnits = calcLpUnits(
        amounts,
        nativeAssetAmount,
        externalAssetAmount
      );

      const poolUnits = lpUnits.add(this.poolUnits);

      return [poolUnits, lpUnits];
    },

    // calculateWithdrawal(
    //   nativeAssetAmount: AssetAmount,
    //   externalAssetAmount: AssetAmount,
    //   lpUnits: Fraction,
    //   wBasisPoints: Fraction,
    //   asymmetry: Fraction
    // ) {
    //   const nativeAssetBalance = amounts.find(
    //     (a) => a.asset.symbol === nativeAssetAmount.asset.symbol
    //   );
    //   const externalAssetBalance = amounts.find(
    //     (a) => a.asset.symbol === externalAssetAmount.asset.symbol
    //   );
    //   if (!nativeAssetBalance || !externalAssetBalance) return null;

    //   const {
    //     lpUnitsLeft,
    //     swapAmount,
    //     withdrawExternalAssetAmount,
    //     withdrawNativeAssetAmount,
    //   } = calculateWithdrawal(
    //     this.poolUnits,
    //     nativeAssetBalance,
    //     externalAssetBalance,
    //     lpUnits,
    //     wBasisPoints,
    //     asymmetry
    //   );
    // },
    // calculateWithdrawal(lpUnits: Fraction, wBasisPoints: Fraction) {
    //   // calculateWithdrawal(this.poolUnits, this.)
    // },
    // {
    //   unitsToClaim = lpUnits / (10000 / wBasisPoints)
    //   withdrawExternalAssetAmount = externalAssetBalance / (poolUnits / unitsToClaim)
    //   withdrawNativeAssetAmount = nativeAssetBalance / (poolUnits / unitsToClaim)

    //   swapAmount = 0
    //   //if asymmetry is positive we need to swap from native to external
    //   if asymmetry > 0
    //     unitsToSwap = (unitsToClaim / (10000 / asymmetry))
    //     swapAmount = nativeAssetBalance / (poolUnits / unitsToSwap)

    //   //if asymmetry is negative we need to swap from external to native
    //   if asymmetry < 0
    //     unitsToSwap = (unitsToClaim / (10000 / asymmetry))
    //     swapAmount = externalAssetBalance / (poolUnits / unitsToSwap)

    //   //if asymmetry is 0 we don't need to swap

    //   lpUnitsLeft = lpUnits - unitsToClaim

    //   return withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount
    // }
  };

  return instance;
}

export function CompositePool(pair1: Pool, pair2: Pool): IPool {
  // The combined asset is the
  const pair1Assets = pair1.amounts.map((a) => a.asset.symbol);
  const pair2Assets = pair2.amounts.map((a) => a.asset.symbol);

  const nativeSymbol = pair1Assets.find((value) => pair2Assets.includes(value));

  if (!nativeSymbol) {
    throw new Error(
      "Cannot create composite pair because pairs do not share a common symbol"
    );
  }

  const amounts = [
    ...pair1.amounts.filter((a) => a.asset.symbol !== nativeSymbol),
    ...pair2.amounts.filter((a) => a.asset.symbol !== nativeSymbol),
  ];

  if (amounts.length !== 2) {
    throw new Error(
      "Cannot create composite pair because pairs do not share a common symbol"
    );
  }

  return {
    amounts: amounts as [IAssetAmount, IAssetAmount],

    priceAsset(asset: Asset) {
      return this.calcSwapResult(AssetAmount(asset, "1"));
    },

    otherAsset(asset: Asset) {
      const otherAsset = amounts.find(
        (amount) => amount.asset.symbol !== asset.symbol
      );
      if (!otherAsset) throw new Error("Asset doesnt exist in pair");
      return otherAsset.asset;
    },

    symbol() {
      return amounts
        .map((a) => a.asset.symbol)
        .sort()
        .join("_");
    },

    contains(...assets: Asset[]) {
      const local = amounts
        .map((a) => a.asset.symbol)
        .sort()
        .join(",");

      const other = assets
        .map((a) => a.symbol)
        .sort()
        .join(",");

      return local === other;
    },

    calcSwapResult(x: AssetAmount) {
      // TODO: possibly use a combined formula
      const [first, second] = pair1.contains(x.asset)
        ? [pair1, pair2]
        : [pair2, pair1];

      const nativeAmount = first.calcSwapResult(x);

      return second.calcSwapResult(nativeAmount);
    },

    calcReverseSwapResult(S: AssetAmount) {
      // TODO: possibly use a combined formula
      const [first, second] = pair1.contains(S.asset)
        ? [pair1, pair2]
        : [pair2, pair1];

      const nativeAmount = first.calcReverseSwapResult(S);

      return second.calcReverseSwapResult(nativeAmount);
    },

    toString() {
      return amounts.map((a) => a.toString()).join(" | ");
    },
  };
}
