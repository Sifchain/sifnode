// Substantially influenced by https://github.com/Uniswap/uniswap-sdk/blob/v2/src/entities/fractions/fraction.ts
/* 
MIT License

Copyright (c) 2020 Noah Zinsmeister

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import invariant from "tiny-invariant";
import JSBI from "jsbi";
import _Decimal from "decimal.js-light";
import _Big from "big.js";
import toFormat from "toformat";

export type BigintIsh = JSBI | bigint | string;

export enum Rounding {
  ROUND_DOWN,
  ROUND_HALF_UP,
  ROUND_UP,
}

export const ZERO = JSBI.BigInt(0);
export const ONE = JSBI.BigInt(1);
export const TWO = JSBI.BigInt(2);
export const THREE = JSBI.BigInt(3);
export const FIVE = JSBI.BigInt(5);
export const TEN = JSBI.BigInt(10);
export const _100 = JSBI.BigInt(100);
export const _997 = JSBI.BigInt(997);
export const _1000 = JSBI.BigInt(1000);

export function parseBigintIsh(bigintIsh: BigintIsh): JSBI {
  return bigintIsh instanceof JSBI
    ? bigintIsh
    : typeof bigintIsh === "bigint"
    ? JSBI.BigInt(bigintIsh.toString())
    : JSBI.BigInt(bigintIsh);
}

const Decimal = toFormat(_Decimal);
const Big = toFormat(_Big);

const toSignificantRounding = {
  [Rounding.ROUND_DOWN]: Decimal.ROUND_DOWN,
  [Rounding.ROUND_HALF_UP]: Decimal.ROUND_HALF_UP,
  [Rounding.ROUND_UP]: Decimal.ROUND_UP,
};

const toFixedRounding = {
  [Rounding.ROUND_DOWN]: 0,
  [Rounding.ROUND_HALF_UP]: 1,
  [Rounding.ROUND_UP]: 3,
};

export class Fraction {
  public readonly numerator: JSBI;
  public readonly denominator: JSBI;

  public constructor(numerator: BigintIsh, denominator: BigintIsh = ONE) {
    this.numerator = parseBigintIsh(numerator);
    this.denominator = parseBigintIsh(denominator);
  }

  // performs floor division
  public get quotient(): JSBI {
    return JSBI.divide(this.numerator, this.denominator);
  }

  // remainder after floor division
  public get remainder(): Fraction {
    return new Fraction(
      JSBI.remainder(this.numerator, this.denominator),
      this.denominator
    );
  }

  public invert(): Fraction {
    return new Fraction(this.denominator, this.numerator);
  }

  public add(other: Fraction | BigintIsh): Fraction {
    const otherParsed =
      other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
    if (JSBI.equal(this.denominator, otherParsed.denominator)) {
      return new Fraction(
        JSBI.add(this.numerator, otherParsed.numerator),
        this.denominator
      );
    }
    return new Fraction(
      JSBI.add(
        JSBI.multiply(this.numerator, otherParsed.denominator),
        JSBI.multiply(otherParsed.numerator, this.denominator)
      ),
      JSBI.multiply(this.denominator, otherParsed.denominator)
    );
  }

  public subtract(other: Fraction | BigintIsh): Fraction {
    const otherParsed =
      other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
    if (JSBI.equal(this.denominator, otherParsed.denominator)) {
      return new Fraction(
        JSBI.subtract(this.numerator, otherParsed.numerator),
        this.denominator
      );
    }
    return new Fraction(
      JSBI.subtract(
        JSBI.multiply(this.numerator, otherParsed.denominator),
        JSBI.multiply(otherParsed.numerator, this.denominator)
      ),
      JSBI.multiply(this.denominator, otherParsed.denominator)
    );
  }

  public lessThan(other: Fraction | BigintIsh): boolean {
    const otherParsed =
      other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
    return JSBI.lessThan(
      JSBI.multiply(this.numerator, otherParsed.denominator),
      JSBI.multiply(otherParsed.numerator, this.denominator)
    );
  }

  public equalTo(other: Fraction | BigintIsh): boolean {
    const otherParsed =
      other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
    return JSBI.equal(
      JSBI.multiply(this.numerator, otherParsed.denominator),
      JSBI.multiply(otherParsed.numerator, this.denominator)
    );
  }

  public greaterThan(other: Fraction | BigintIsh): boolean {
    const otherParsed =
      other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
    return JSBI.greaterThan(
      JSBI.multiply(this.numerator, otherParsed.denominator),
      JSBI.multiply(otherParsed.numerator, this.denominator)
    );
  }

  public multiply(other: Fraction | BigintIsh): Fraction {
    const otherParsed =
      other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
    return new Fraction(
      JSBI.multiply(this.numerator, otherParsed.numerator),
      JSBI.multiply(this.denominator, otherParsed.denominator)
    );
  }

  public divide(other: Fraction | BigintIsh): Fraction {
    const otherParsed =
      other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
    return new Fraction(
      JSBI.multiply(this.numerator, otherParsed.denominator),
      JSBI.multiply(this.denominator, otherParsed.numerator)
    );
  }

  public toSignificant(
    significantDigits: number,
    format: object = { groupSeparator: "" },
    rounding: Rounding = Rounding.ROUND_HALF_UP
  ): string {
    invariant(
      Number.isInteger(significantDigits),
      `${significantDigits} is not an integer.`
    );
    invariant(significantDigits > 0, `${significantDigits} is not positive.`);

    Decimal.set({
      precision: significantDigits + 1,
      rounding: toSignificantRounding[rounding],
    });
    const quotient = new Decimal(this.numerator.toString())
      .div(this.denominator.toString())
      .toSignificantDigits(significantDigits);
    return quotient.toFormat(quotient.decimalPlaces(), format);
  }

  public toFixed(
    decimalPlaces: number,
    format: object = { groupSeparator: "" },
    rounding: Rounding = Rounding.ROUND_HALF_UP
  ): string {
    invariant(
      Number.isInteger(decimalPlaces),
      `${decimalPlaces} is not an integer.`
    );
    invariant(decimalPlaces >= 0, `${decimalPlaces} is negative.`);

    Big.DP = decimalPlaces;
    Big.RM = toFixedRounding[rounding];
    return new Big(this.numerator.toString())
      .div(this.denominator.toString())
      .toFormat(decimalPlaces, format);
  }
}
