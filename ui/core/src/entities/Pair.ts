import { Asset } from "./Asset";
import { IAssetAmount } from "./AssetAmount";

export type Pair = ReturnType<typeof Pair>;

export function Pair(a: IAssetAmount, b: IAssetAmount) {
  const amounts: [IAssetAmount, IAssetAmount] = [a, b];

  return {
    amounts,

    otherAsset(asset: Asset) {
      const otherAsset = amounts.find(
        (amount) => amount.symbol !== asset.symbol,
      );
      if (!otherAsset) throw new Error("Asset doesnt exist in pair");
      return otherAsset;
    },

    symbol() {
      return amounts
        .map((a) => a.symbol)
        .sort()
        .join("_");
    },

    contains(...assets: Asset[]) {
      const local = amounts.map((a) => a.symbol);

      const other = assets.map((a) => a.symbol);

      return !!local.find((s) => other.includes(s));
    },

    getAmount(asset: Asset | string) {
      const assetSymbol = typeof asset === "string" ? asset : asset.symbol;
      const found = this.amounts.find((amount) => {
        return amount.symbol === assetSymbol;
      });
      if (!found) throw new Error(`Asset ${assetSymbol} doesnt exist in pair`);
      return found;
    },

    toString() {
      return amounts.map((a) => a.toString()).join(" | ");
    },
  };
}
