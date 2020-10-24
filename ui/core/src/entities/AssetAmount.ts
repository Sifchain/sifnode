import { Asset } from "./Asset";
import invariant from "tiny-invariant";
import _Big from "big.js";
import toFormat from "toformat";

import {
  BigintIsh,
  Fraction,
  parseBigintIsh,
  Rounding,
  TEN,
} from "./fraction/Fraction";
import JSBI from "jsbi";
import B from "./utils/B";

const Big = toFormat(_Big);

export class AssetAmountFraction extends Fraction {
  constructor(public asset: Asset, public amount: BigintIsh) {
    super(
      parseBigintIsh(amount),
      JSBI.exponentiate(TEN, JSBI.BigInt(asset.decimals))
    );
  }

  public toSignificant(
    significantDigits = 6,
    format?: object,
    rounding: Rounding = Rounding.ROUND_DOWN
  ): string {
    return super.toSignificant(significantDigits, format, rounding);
  }

  public toFixed(
    decimalPlaces = this.asset.decimals,
    format?: object,
    rounding: Rounding = Rounding.ROUND_DOWN
  ): string {
    invariant(decimalPlaces <= this.asset.decimals, "DECIMALS");
    return super.toFixed(decimalPlaces, format, rounding);
  }

  public toExact(format: object = { groupSeparator: "" }): string {
    Big.DP = this.asset.decimals;
    return new Big(this.numerator.toString())
      .div(this.denominator.toString())
      .toFormat(format);
  }

  public toString() {
    return `${this.asset.symbol} ${this.toFixed()}`;
  }
}

export type AssetAmount = AssetAmountFraction;

// Conveniance method for initializing a balance with a number
// AssetAmountN(ETH, 10)
// If amount is in JSBI then the amount this creates will be in base units (ie Wei)
// Otherwise the amount will be in natural units
export function AssetAmountN(
  asset: Asset,
  amount: string | number | JSBI
): AssetAmount {
  const jsbiAmount =
    amount instanceof JSBI ? amount : B(amount, asset.decimals);
  return new AssetAmountFraction(asset, jsbiAmount);
}
