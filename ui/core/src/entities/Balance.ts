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
const Big = toFormat(_Big);

export class Balance extends Fraction {
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
  static create(asset: Asset, amount: BigintIsh): Balance {
    return new Balance(asset, amount);
  }
}

export type Balances = Balance[];

// export type AssetBalancesByAddress = {
//   [address: string]: Balance | undefined;
// };
