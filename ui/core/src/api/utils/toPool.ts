import { RWN } from "../../constants";
import { Asset, AssetAmount, Fraction, Pool } from "../../entities";
import { RawPool } from "./x/clp";

export function toPool(poolData: RawPool): Pool {
  console.log({ poolData });
  const externalAssetTicker = poolData.external_asset.ticker;

  return Pool(
    AssetAmount(RWN, poolData.native_asset_balance),
    AssetAmount(
      Asset.get(externalAssetTicker),
      poolData.external_asset_balance
    ),
    new Fraction(poolData.pool_units)
  );
}
