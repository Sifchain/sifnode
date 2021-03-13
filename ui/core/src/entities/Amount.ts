import JSBI from "jsbi";
import { Fraction, IFraction } from "./fraction/Fraction";
import Big from "big.js";

export type IAmount = {
  // for use by display lib and in testing
  toBigInt(): JSBI;
  toString(): string;

  // for use elsewhere
  add(other: IAmount | string): IAmount;
  divide(other: IAmount | string): IAmount;
  equalTo(other: IAmount | string): boolean;
  greaterThan(other: IAmount | string): boolean;
  greaterThanOrEqual(other: IAmount | string): boolean;
  lessThan(other: IAmount | string): boolean;
  lessThanOrEqual(other: IAmount | string): boolean;
  multiply(other: IAmount | string): IAmount;
  sqrt(): IAmount;
  subtract(other: IAmount | string): IAmount;
};

// exported ONLY to be shared with AssetAmount!
export type _ExposeInternal<T extends IAmount> = T & {
  _toInternal(): IFraction;
};

function toFraction(a: IAmount | string): IFraction | string {
  type _IAmount = _ExposeInternal<IAmount>;
  if (typeof a === "string") return a;
  return (a as _IAmount)._toInternal();
}

function toBig(fraction: Fraction) {
  return Big(fraction.toFixed(0));
}

function toAmount(a: IFraction) {
  return Amount(a.quotient);
}

export function Amount(source: JSBI | bigint | string | IAmount): IAmount {
  type _IAmount = _ExposeInternal<IAmount>;

  if (
    !(source instanceof JSBI) &&
    typeof source !== "bigint" &&
    typeof source !== "string"
  ) {
    return source;
  }

  const fraction = new Fraction(source);
  const instance: _IAmount = {
    toBigInt() {
      return fraction.quotient;
    },

    toString() {
      return fraction.toFixed(0);
    },

    add(other) {
      return toAmount(fraction.add(toFraction(other)));
    },

    divide(other) {
      return toAmount(fraction.divide(toFraction(other)));
    },

    equalTo(other) {
      return fraction.equalTo(toFraction(other));
    },

    greaterThan(other) {
      return fraction.greaterThan(toFraction(other));
    },

    greaterThanOrEqual(other) {
      return fraction.greaterThanOrEqual(toFraction(other));
    },

    lessThan(other) {
      return fraction.lessThan(toFraction(other));
    },

    lessThanOrEqual(other) {
      return fraction.lessThanOrEqual(toFraction(other));
    },

    multiply(other) {
      return toAmount(fraction.multiply(toFraction(other)));
    },

    subtract(other) {
      return toAmount(fraction.subtract(toFraction(other)));
    },

    sqrt() {
      // TODO: test against rounding errors
      const big = toBig(fraction);
      const string = toFraction(big.sqrt().toFixed(0)) as string;
      return Amount(string);
    },

    _toInternal() {
      return fraction;
    },
  };

  return instance;
}
