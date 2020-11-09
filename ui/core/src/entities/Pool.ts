import { Asset } from "./Asset";
import { AssetAmount } from "./AssetAmount";
import { Pair } from "./Pair";
import Big from "big.js";
import JSBI from "jsbi";

export type Pool = ReturnType<typeof Pool>;
export type IPool = Omit<Pool, "poolUnits" | "calculatePoolUnits">;

export function Pool(
  a: AssetAmount,
  b: AssetAmount,
  poolUnits: JSBI = JSBI.BigInt("0")
) {
  const pair = Pair(a, b);
  const amounts = pair.amounts;
  return {
    amounts,

    otherAsset: pair.otherAsset,
    symbol: pair.symbol,
    contains: pair.contains,
    toString: pair.toString,
    poolUnits,
    priceAsset(asset: Asset) {
      return this.calcSwapResult(AssetAmount(asset, "1"));
    },

    // https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
    // Formula: swapAmount = (x * X * Y) / (x + X) ^ 2
    calcSwapResult(x: AssetAmount) {
      const X = amounts.find((a) => a.asset.symbol === x.asset.symbol);
      if (!X) throw new Error("Sent amount does not exist in this pair");
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
      if (!Ya) throw new Error("Sent amount does not exist in this pair");
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

    // https://github.com/Sifchain/sifnode/blob/develop/docs/1.Liquidity%20Pools%20Architecture.md
    // poolerUnits = ((R + A) * (r * A + R * a))/(4 * R * A)
    calculatePoolUnits(
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
      const lpUnits = numerator.divide(denominator);
      const poolUnits = JSBI.add(this.poolUnits, lpUnits.quotient);
      return [poolUnits, lpUnits];
    },
  };
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
