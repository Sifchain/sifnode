import { Asset } from "./Asset";
import { AssetAmount } from "./AssetAmount";
import { Fraction } from "./fraction/Fraction";

export type Pair = {
  amounts: [AssetAmount, AssetAmount];
  priceA: () => Fraction;
  priceB: () => Fraction;
};

export function Pair(a: AssetAmount, b: AssetAmount) {
  const amounts = [a, b];

  return {
    amounts,
    priceA() {
      return AssetAmount(b.asset, b.divide(a).toFixed(b.asset.decimals));
    },

    priceB() {
      return AssetAmount(a.asset, a.divide(b).toFixed(a.asset.decimals));
    },

    contains(...assets: Asset[]) {
      const me = amounts
        .map((a) => a.asset.symbol)
        .sort()
        .join(",");

      const you = assets
        .map((a) => a.symbol)
        .sort()
        .join(",");

      return me == you;
    },
  };
}
