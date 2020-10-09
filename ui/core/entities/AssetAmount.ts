import { Asset } from "./Asset";

export type AssetAmount = {
  asset: Asset;
  amount: BigInt;
};

export type AssetBalancesByAddress = {
  [address: string]: AssetAmount | undefined;
};

export function createAssetAmount(asset: Asset, amount: BigInt): AssetAmount {
  return {
    asset,
    amount,
  };
}
