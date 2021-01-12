import { Asset, AssetAmount, Coin, Fraction, Pool } from "../../../entities";
import { RawPool } from "./x/clp";

export const toPool = (nativeAsset: Coin) => (poolData: RawPool): Pool => {
  const externalAssetSymbol = poolData.external_asset.symbol;

  return Pool(
    AssetAmount(nativeAsset, poolData.native_asset_balance, {
      inBaseUnit: true,
    }),
    AssetAmount(
      Asset.get(externalAssetSymbol),
      poolData.external_asset_balance,
      { inBaseUnit: true }
    ),
    new Fraction(poolData.pool_units)
  );
};
