import JSBI from "jsbi";
import { Fraction, IFraction } from "./fraction/Fraction";
import Big from "big.js";

export type IAmount = {
  // for use by display lib and in testing
  toBigInt(): JSBI;
  toString(detailed?: boolean): string;

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
  _fromInternal(fraction: IFraction): IAmount;
};

function toFraction(a: string): string;
function toFraction(a: IAmount | string): IFraction;
function toFraction(a: IAmount | string): IFraction | string {
  type _IAmount = _ExposeInternal<IAmount>;
  if (typeof a === "string") return a;
  return (a as _IAmount)._toInternal();
}

function toBig(fraction: Fraction) {
  return Big(fraction.toFixed(0));
}

function toAmount(a: IFraction) {
  type _IAmount = _ExposeInternal<IAmount>;
  return (Amount("0") as _IAmount)._fromInternal(a);
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

  let fraction = new Fraction(source);
  const instance: _IAmount = {
    toBigInt() {
      return getQuotientWithBankersRounding(fraction);
    },

    toString(detailed?: boolean) {
      return fraction.toFixed(detailed ? 18 : 0);
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

    _fromInternal(_fraction: IFraction) {
      fraction = _fraction;
      return instance;
    },

    _toInternal() {
      return fraction;
    },
  };

  return instance;
}

// quotient needs to use bankers rounding so we follow this example for bankers rounding in BigInt and apply to JSBI
//https://stackoverflow.com/questions/53752370/ecmascript-bigint-round-to-even
function getQuotientWithBankersRounding(fraction: IFraction): JSBI {
  const a = fraction.numerator;
  const b = fraction.denominator;

  const aAbs = JSBI.greaterThan(a, JSBI.BigInt("0"))
    ? a
    : JSBI.multiply(JSBI.BigInt("-1"), a);

  const bAbs = JSBI.greaterThan(b, JSBI.BigInt("0"))
    ? b
    : JSBI.multiply(JSBI.BigInt("-1"), b);

  let result = JSBI.divide(aAbs, bAbs);

  const rem = JSBI.remainder(aAbs, bAbs);

  if (JSBI.greaterThan(JSBI.multiply(rem, JSBI.BigInt("2")), bAbs)) {
    result = JSBI.add(result, JSBI.BigInt("1"));
  } else if (JSBI.equal(JSBI.multiply(rem, JSBI.BigInt("2")), bAbs)) {
    if (
      JSBI.equal(JSBI.remainder(result, JSBI.BigInt("2")), JSBI.BigInt("1"))
    ) {
      result = JSBI.add(result, JSBI.BigInt("1"));
    }
  }

  if (
    JSBI.greaterThan(a, JSBI.BigInt("0")) !==
    JSBI.greaterThan(b, JSBI.BigInt("0"))
  ) {
    return JSBI.multiply(JSBI.BigInt("-1"), result);
  } else {
    return result;
  }
}
