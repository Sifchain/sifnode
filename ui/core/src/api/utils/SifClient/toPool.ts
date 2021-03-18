import { Asset, AssetAmount, Fraction, Pool } from "../../../entities";
import { RawPool } from "./x/clp";

function getAssetOrNull(symbol: string): Asset | null {
  try {
    return Asset.get(symbol);
  } catch (err) {
    return null;
  }
}
export const toPool = (nativeAsset: Asset) => (
  poolData: RawPool,
): Pool | null => {
  const externalAssetSymbol = poolData.external_asset.symbol;
  const externalAsset = getAssetOrNull(externalAssetSymbol);

  // If we are not configured to handle this external asset
  // the pool is invalid so we ignore it
  if (!externalAsset) return null;

  return Pool(
    AssetAmount(nativeAsset, poolData.native_asset_balance, {
      inBaseUnit: true,
    }),
    AssetAmount(externalAsset, poolData.external_asset_balance, {
      inBaseUnit: true,
    }),
    new Fraction(poolData.pool_units),
  );
};
