import { Asset } from "./Asset";
import { AssetAmount } from "./AssetAmount";
import { Fraction } from "./fraction/Fraction";
import Big from "big.js";
export type Pair = ReturnType<typeof Pair>;

export function Pair(a: AssetAmount, b: AssetAmount) {
  const amounts = [a, b];

  return {
    amounts,

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
      const local = amounts.map((a) => a.asset.symbol);

      const other = assets.map((a) => a.symbol);

      return !!local.find((s) => other.includes(s));
    },

    // https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
    // Formula: swapAmount = (x * X * Y) / (x + X) ^ 2

    calcSwapResult(x: AssetAmount) {
      const X = amounts.find((a) => a.asset.symbol === x.asset.symbol);
      if (!X) throw new Error("Sent amount does not exist in this pair");
      const Y = amounts.find((a) => a.asset.symbol !== x.asset.symbol);
      if (!Y) throw new Error("Pair does not have an opposite asset."); // For Typescript's sake will probably never happen

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
      if (!Ya) throw new Error("Sent amount does not exist in this pair");
      const Xa = amounts.find((a) => a.asset.symbol !== Sa.asset.symbol);
      if (!Xa) throw new Error("Pair does not have an opposite asset."); // For Typescript's sake will probably never happen
      const otherAsset = this.otherAsset(Sa.asset);
      if (Sa.equalTo("0")) {
        return AssetAmount(otherAsset, "0");
      }

      // Need to use Big.js for sqrt calculation
      // Ok to accept a little precision loss as reverse swap amount can be rough
      const S = Big(Sa.toFixed());
      const X = Big(Xa.toFixed());
      const Y = Big(Ya.toFixed());

      const term1 = Big(-2)
        .times(X)
        .times(S);

      const term2 = X.times(Y);
      const underRoot = Y.times(Y.minus(S.times(4)));

      const term3 = X.times(underRoot.sqrt());

      const numerator = term1.plus(term2).minus(term3);
      const denominator = S.times(2);

      const x = numerator.div(denominator);

      return AssetAmount(otherAsset, x.toFixed());
    },
    toString() {
      return amounts.map((a) => a.toString()).join(" | ");
    },
  };
}

export function CompositePair(pair1: Pair, pair2: Pair): Pair {
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
    amounts,
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
      const [first, second] = pair1.contains(x.asset)
        ? [pair1, pair2]
        : [pair2, pair1];

      const nativeAmount = first.calcSwapResult(x);

      return second.calcSwapResult(nativeAmount);
    },

    calcReverseSwapResult(S: AssetAmount) {
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
