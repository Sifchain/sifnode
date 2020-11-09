import { Asset } from "./Asset";
import { AssetAmount } from "./AssetAmount";

export type Pair = ReturnType<typeof Pair>;

export function Pair(a: AssetAmount, b: AssetAmount) {
  const amounts: [AssetAmount, AssetAmount] = [a, b];

  return {
    amounts,

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

    toString() {
      return amounts.map((a) => a.toString()).join(" | ");
    },
  };
}
