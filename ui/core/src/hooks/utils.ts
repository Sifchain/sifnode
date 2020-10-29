import { computed, Ref } from "@vue/reactivity";
import { Asset, AssetAmount, Pair } from "../entities";

export function assetPriceMessage(asset: Asset | null, pair: Pair | null) {
  if (!asset || !pair) return "";
  const assetPrice = pair.priceAsset(asset);

  if (!assetPrice || (assetPrice && assetPrice.equalTo("0"))) return "";

  return `${assetPrice.toFormatted()} per ${asset?.symbol.toUpperCase()}`;
}

export function useBalances(balances: Ref<AssetAmount[]>) {
  return computed(() => {
    const map = new Map<string, AssetAmount>();

    for (const item of balances.value) {
      map.set(item.asset.symbol, item);
    }
    return map;
  });
}
