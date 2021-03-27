import JSBI from "jsbi";
import { Fraction, IFraction } from "./fraction/Fraction";
import Big from "big.js";
import { isAssetAmount } from "./AssetAmount";
import { decimalShift, floorDecimal } from "../utils/decimalShift";

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
  powerInt(other: IAmount | string): IAmount;
  expInt(other: IAmount | string): IAmount;
  sqrt(): IAmount;
  subtract(other: IAmount | string): IAmount;
};

export function Amount(
  source: JSBI | bigint | string | IAmount,
): Readonly<IAmount> {
  type _IAmount = _ExposeInternal<IAmount>;

  if (typeof source === "string" && source.match(/^[+-]?\s?(\d+)?\.\d+$/)) {
    return getAmountFromDecimal(source);
  }

  if (
    !(source instanceof JSBI) &&
    typeof source !== "bigint" &&
    typeof source !== "string"
  ) {
    if (isAssetAmount(source)) {
      return source.amount;
    }
    return source;
  }

  let fraction = new Fraction(source);
  const instance: _IAmount = {
    // We only loose precision and round when we move to BigInt for display
    toBigInt() {
      return getQuotientWithBankersRounding(fraction);
    },

    toString(detailed: boolean = true) {
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

    powerInt(other) {
      const exp = parseInt(Amount(other).toBigInt().toString());

      if (exp === 0) {
        if (fraction.greaterThan("0")) return Amount("1");
        return Amount("-1");
      }

      let fr = fraction;
      for (let i = Math.abs(exp); i > 1; --i) {
        fr = fr.multiply(fraction);
      }
      if (exp >= 0) {
        return toAmount(fr);
      }

      return toAmount(new Fraction("1").divide(fr));
    },

    expInt(other) {
      // TODO: use decimalShift for speed - but this might loose precision?
      return instance.multiply(Amount("10").powerInt(other));
    },

    sqrt() {
      // TODO: test against rounding errors
      const big = toBig(fraction);
      const string = toFraction(
        big.sqrt().times("100000000000000000000000").toFixed(0),
      ) as string;
      return Amount(string).divide("100000000000000000000000");
    },

    // Internal methods need to be exposed here
    // so they can be used by another Amount in
    // toFraction and toAmount
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

function getAmountFromDecimal(decimal: string): IAmount {
  return Amount(floorDecimal(decimalShift(decimal, 18))).divide(
    "1000000000000000000",
  );
}

// exported ONLY to be shared with AssetAmount!
export type _ExposeInternal<T extends IAmount> = T & {
  // Private method to expose internal representation
  _toInternal(): IFraction;

  // Private method to populate IAmount value from internal representation
  _fromInternal(fraction: IFraction): IAmount;
};

// Helper for extracting a fraction out of an amount.
// This uses a private API and should not be exposed
// outside of Amount
function toFraction(a: string): string;
function toFraction(a: IAmount | string): IFraction;
function toFraction(a: IAmount | string): IFraction | string {
  type _IAmount = _ExposeInternal<IAmount>;
  if (typeof a === "string") {
    return a.indexOf(".") < 0 ? a : (Amount(a) as _IAmount)._toInternal();
  }
  return (a as _IAmount)._toInternal();
}

// Internal helper convert to Big.js for calculating sqrts
// NOTE this looses precision to 1e24
function toBig(fraction: Fraction) {
  return Big(fraction.toFixed(24));
}

// Helper for converting a fraction to an amount.
// This uses a private API and should not be exposed
// outside of Amount
function toAmount(a: IFraction) {
  type _IAmount = _ExposeInternal<IAmount>;
  return (Amount("0") as _IAmount)._fromInternal(a);
}
