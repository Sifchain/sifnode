import { Asset } from "./Asset";
import { IFraction } from "./fraction/Fraction";

export function LiquidityProvider(
  asset: Asset,
  units: IFraction,
  address: string,
  nativeAmount: IFraction,
  externalAmount: IFraction
) {
  return { asset, units, address, nativeAmount, externalAmount };
}
export type LiquidityProvider = {
  asset: Asset;
  units: IFraction;
  address: string;
  nativeAmount: IFraction;
  externalAmount: IFraction;
};
