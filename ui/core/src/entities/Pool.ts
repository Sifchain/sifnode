import { Asset } from "./Asset";
import { AssetAmount, IAssetAmount } from "./AssetAmount";
import { Pair } from "./Pair";
import Big from "big.js";
import { Fraction } from "./fraction/Fraction";
import {
  calculatePoolUnits,
  calculatePriceImpact,
  calculateProviderFee,
  calculateReverseSwapResult,
  calculateSwapResult,
} from "./formulae";

export type Pool = ReturnType<typeof Pool>;
export type IPool = Omit<Pool, "poolUnits" | "calculatePoolUnits">;

export function Pool(
  a: AssetAmount, // native asset
  b: AssetAmount, // external asset
  poolUnits?: Fraction
) {
  const pair = Pair(a, b);
  const amounts: [IAssetAmount, IAssetAmount] = pair.amounts;

  return {
    amounts,
    otherAsset: pair.otherAsset,
    symbol: pair.symbol,
    contains: pair.contains,
    toString: pair.toString,
    getAmount: pair.getAmount,
    poolUnits:
      poolUnits ||
      calculatePoolUnits(
        a,
        b,
        new Fraction("0"),
        new Fraction("0"),
        new Fraction("0")
      ),
    priceAsset(asset: Asset) {
      return this.calcSwapResult(AssetAmount(asset, "1"));
    },

    calcProviderFee(x: AssetAmount) {
      const X = amounts.find((a) => a.asset.symbol === x.asset.symbol);
      if (!X)
        throw new Error(
          `Sent amount with symbol ${
            x.asset.symbol
          } does not exist in this pair: ${this.toString()}`
        );
      const Y = amounts.find((a) => a.asset.symbol !== x.asset.symbol);
      if (!Y) throw new Error("Pool does not have an opposite asset."); // For Typescript's sake will probably never happen
      const providerFee = calculateProviderFee(x, X, Y);
      return AssetAmount(this.otherAsset(x.asset), providerFee);
    },

    calcPriceImpact(x: AssetAmount) {
      const X = amounts.find((a) => a.asset.symbol === x.asset.symbol);
      if (!X)
        throw new Error(
          `Sent amount with symbol ${
            x.asset.symbol
          } does not exist in this pair: ${this.toString()}`
        );
      return calculatePriceImpact(x, X).multiply("100");
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
      const swapAmount = calculateSwapResult(x, X, Y);
      return AssetAmount(this.otherAsset(x.asset), swapAmount);
    },

    calcReverseSwapResult(Sa: AssetAmount): IAssetAmount {
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

    calculatePoolUnits(
      nativeAssetAmount: AssetAmount,
      externalAssetAmount: AssetAmount
    ) {
      const [nativeBalanceBefore, externalBalanceBefore] = amounts;

      // Calculate current units created by this potential liquidity provision
      const lpUnits = calculatePoolUnits(
        nativeAssetAmount,
        externalAssetAmount,
        nativeBalanceBefore,
        externalBalanceBefore,
        this.poolUnits
      );
      const newTotalPoolUnits = lpUnits.add(this.poolUnits);

      return [newTotalPoolUnits, lpUnits];
    },
  };
}

export function CompositePool(pair1: IPool, pair2: IPool): IPool {
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

    getAmount: (asset: Asset | string) => {
      if (Asset.get(asset).symbol === nativeSymbol) {
        throw new Error(`Asset ${nativeSymbol} doesnt exist in pair`);
      }

      // quicker to try catch than contains
      try {
        return pair1.getAmount(asset);
      } catch (err) {}

      return pair2.getAmount(asset);
    },

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
      const local = amounts.map((a) => a.asset.symbol).sort();
      const other = assets.map((a) => a.symbol).sort();
      return !!local.find((s) => other.includes(s));
    },

    calcProviderFee(x: AssetAmount) {
      const [first, second] = pair1.contains(x.asset)
        ? [pair1, pair2]
        : [pair2, pair1];
      const firstSwapFee = first.calcProviderFee(x);
      const firstSwapOutput = first.calcSwapResult(x);
      const secondSwapFee = second.calcProviderFee(firstSwapOutput);
      const firstSwapFeeInOutputAsset = second.calcSwapResult(firstSwapFee);
      return AssetAmount(
        second.otherAsset(firstSwapFee.asset),
        firstSwapFeeInOutputAsset.add(secondSwapFee)
      );
    },

    calcPriceImpact(x: AssetAmount) {
      const [first, second] = pair1.contains(x.asset)
        ? [pair1, pair2]
        : [pair2, pair1];
      const firstPoolImpact = first.calcPriceImpact(x);
      const r = first.calcSwapResult(x);
      const secondPoolImpact = second.calcPriceImpact(r);
      return firstPoolImpact.add(secondPoolImpact);
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
