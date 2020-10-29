import { Asset } from "./Asset";
import { AssetAmount } from "./AssetAmount";

export type Pair = ReturnType<typeof Pair>;

const hasAsset = (asset: Asset) => (amount: AssetAmount) => {
  return amount.asset === asset;
};

export function Pair(a: AssetAmount, b: AssetAmount) {
  const amounts = [a, b];

  return {
    amounts,
    priceA() {
      const asset = b.asset;
      if (a.equalTo("0")) return null;
      return AssetAmount(asset, b.divide(a).toFixed(asset.decimals));
    },

    priceB() {
      const asset = a.asset;
      if (b.equalTo("0")) return null;
      return AssetAmount(asset, a.divide(b).toFixed(asset.decimals));
    },

    priceAsset(asset: Asset) {
      if (a.asset === asset) {
        return this.priceA();
      }

      if (b.asset === asset) {
        return this.priceB();
      }

      throw new Error(`Asset not ${asset.symbol} found in pair.`);
    },

    hasAsset(asset: Asset) {
      return amounts.filter((amount) => amount.asset === asset).length > 0;
    },

    oppositeAsset(asset: Asset) {
      if (!hasAsset(asset)) {
        return null;
      }
      return amounts
        .map((amount) => amount.asset)
        .filter((a) => a !== asset)[0];
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
  };
}
