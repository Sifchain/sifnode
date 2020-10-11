import { Asset } from './Asset';
import invariant from 'tiny-invariant';
import _Big from 'big.js';
import toFormat from 'toformat';

import {
  BigintIsh,
  Fraction,
  parseBigintIsh,
  Rounding,
  TEN,
} from './fraction/Fraction';
import JSBI from 'jsbi';
const Big = toFormat(_Big);

export class AssetAmount extends Fraction {
  constructor(public asset: Asset, public amount: BigintIsh) {
    super(
      parseBigintIsh(amount),
      JSBI.exponentiate(TEN, JSBI.BigInt(asset.decimals))
    );
  }
  public toSignificant(
    significantDigits: number = 6,
    format?: object,
    rounding: Rounding = Rounding.ROUND_DOWN
  ): string {
    return super.toSignificant(significantDigits, format, rounding);
  }

  public toFixed(
    decimalPlaces: number = this.asset.decimals,
    format?: object,
    rounding: Rounding = Rounding.ROUND_DOWN
  ): string {
    invariant(decimalPlaces <= this.asset.decimals, 'DECIMALS');
    return super.toFixed(decimalPlaces, format, rounding);
  }

  public toExact(format: object = { groupSeparator: '' }): string {
    Big.DP = this.asset.decimals;
    return new Big(this.numerator.toString())
      .div(this.denominator.toString())
      .toFormat(format);
  }
  static create(asset: Asset, amount: BigintIsh): AssetAmount {
    return new AssetAmount(asset, amount);
  }
}

export type AssetBalancesByAddress = {
  [address: string]: AssetAmount | undefined;
};
