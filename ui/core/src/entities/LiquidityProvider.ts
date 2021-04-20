import { Asset } from "./Asset";
import { IAmount } from "./Amount";

export function LiquidityProvider(
  asset: Asset,
  units: IAmount,
  address: string,
  nativeAmount: IAmount,
  externalAmount: IAmount,
) {
  return { asset, units, address, nativeAmount, externalAmount };
}
export type LiquidityProvider = {
  asset: Asset;
  units: IAmount;
  address: string;
  nativeAmount: IAmount;
  externalAmount: IAmount;
};
