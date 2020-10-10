import { Asset } from "./Asset";
import JSBI from "jsbi";
export type AssetAmount = {
  asset: Asset;
  amount: JSBI;
};

export type AssetBalancesByAddress = {
  [address: string]: AssetAmount | undefined;
};

export function createAssetAmount(asset: Asset, amount: JSBI): AssetAmount {
  return {
    asset,
    amount,
  };
}
