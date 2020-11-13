import { Asset } from "./Asset";
import { Fraction } from "./fraction/Fraction";

export function LiquidityProvider(
  asset: Asset,
  units: Fraction,
  address: string
): LiquidityProvider {
  return { asset, units, address };
}
export type LiquidityProvider = {
  asset: Asset;
  units: Fraction;
  address: string;
};
