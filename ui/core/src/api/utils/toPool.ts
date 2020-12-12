import { ROWAN } from "../../constants";
import { Asset, AssetAmount, Fraction, Pool } from "../../entities";
import { RawPool } from "./x/clp";

export function toPool(poolData: RawPool): Pool {
  const externalAssetTicker = poolData.external_asset.ticker;

  return Pool(
    AssetAmount(ROWAN, poolData.native_asset_balance),
    AssetAmount(
      Asset.get(externalAssetTicker),
      poolData.external_asset_balance
    ),
    new Fraction(poolData.pool_units)
  );
}
