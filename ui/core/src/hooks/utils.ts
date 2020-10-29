import { Asset, Pair } from "../entities";

export function assetPriceMessage(asset: Asset | null, pair: Pair | null) {
  if (!asset || !pair) return "";
  const assetPrice = pair.priceAsset(asset);

  if (!assetPrice || (assetPrice && assetPrice.equalTo("0"))) return "";

  return `${assetPrice.toFormatted()} per ${asset?.symbol.toUpperCase()}`;
}
