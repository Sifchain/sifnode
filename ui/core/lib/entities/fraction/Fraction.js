import invariant from "tiny-invariant";
import JSBI from "jsbi";
import _Decimal from "decimal.js-light";
import _Big from "big.js";
import toFormat from "toformat";
export var Rounding;
(function (Rounding) {
    Rounding[Rounding["ROUND_DOWN"] = 0] = "ROUND_DOWN";
    Rounding[Rounding["ROUND_HALF_UP"] = 1] = "ROUND_HALF_UP";
    Rounding[Rounding["ROUND_UP"] = 2] = "ROUND_UP";
})(Rounding || (Rounding = {}));
export const ZERO = JSBI.BigInt(0);
export const ONE = JSBI.BigInt(1);
export const TWO = JSBI.BigInt(2);
export const THREE = JSBI.BigInt(3);
export const FIVE = JSBI.BigInt(5);
export const TEN = JSBI.BigInt(10);
export const _100 = JSBI.BigInt(100);
export const _997 = JSBI.BigInt(997);
export const _1000 = JSBI.BigInt(1000);
export function parseBigintIsh(bigintIsh) {
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
    constructor(numerator, denominator = ONE) {
        this.numerator = parseBigintIsh(numerator);
        this.denominator = parseBigintIsh(denominator);
    }
    // performs floor division
    get quotient() {
        return JSBI.divide(this.numerator, this.denominator);
    }
    // remainder after floor division
    get remainder() {
        return new Fraction(JSBI.remainder(this.numerator, this.denominator), this.denominator);
    }
    invert() {
        return new Fraction(this.denominator, this.numerator);
    }
    add(other) {
        const otherParsed = other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
        if (JSBI.equal(this.denominator, otherParsed.denominator)) {
            return new Fraction(JSBI.add(this.numerator, otherParsed.numerator), this.denominator);
        }
        return new Fraction(JSBI.add(JSBI.multiply(this.numerator, otherParsed.denominator), JSBI.multiply(otherParsed.numerator, this.denominator)), JSBI.multiply(this.denominator, otherParsed.denominator));
    }
    subtract(other) {
        const otherParsed = other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
        if (JSBI.equal(this.denominator, otherParsed.denominator)) {
            return new Fraction(JSBI.subtract(this.numerator, otherParsed.numerator), this.denominator);
        }
        return new Fraction(JSBI.subtract(JSBI.multiply(this.numerator, otherParsed.denominator), JSBI.multiply(otherParsed.numerator, this.denominator)), JSBI.multiply(this.denominator, otherParsed.denominator));
    }
    lessThan(other) {
        const otherParsed = other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
        return JSBI.lessThan(JSBI.multiply(this.numerator, otherParsed.denominator), JSBI.multiply(otherParsed.numerator, this.denominator));
    }
    equalTo(other) {
        const otherParsed = other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
        return JSBI.equal(JSBI.multiply(this.numerator, otherParsed.denominator), JSBI.multiply(otherParsed.numerator, this.denominator));
    }
    greaterThan(other) {
        const otherParsed = other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
        return JSBI.greaterThan(JSBI.multiply(this.numerator, otherParsed.denominator), JSBI.multiply(otherParsed.numerator, this.denominator));
    }
    multiply(other) {
        const otherParsed = other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
        return new Fraction(JSBI.multiply(this.numerator, otherParsed.numerator), JSBI.multiply(this.denominator, otherParsed.denominator));
    }
    divide(other) {
        const otherParsed = other instanceof Fraction ? other : new Fraction(parseBigintIsh(other));
        return new Fraction(JSBI.multiply(this.numerator, otherParsed.denominator), JSBI.multiply(this.denominator, otherParsed.numerator));
    }
    toSignificant(significantDigits, format = { groupSeparator: "" }, rounding = Rounding.ROUND_HALF_UP) {
        invariant(Number.isInteger(significantDigits), `${significantDigits} is not an integer.`);
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
    toFixed(decimalPlaces, format = { groupSeparator: "" }, rounding = Rounding.ROUND_HALF_UP) {
        invariant(Number.isInteger(decimalPlaces), `${decimalPlaces} is not an integer.`);
        invariant(decimalPlaces >= 0, `${decimalPlaces} is negative.`);
        Big.DP = decimalPlaces;
        Big.RM = toFixedRounding[rounding];
        return new Big(this.numerator.toString())
            .div(this.denominator.toString())
            .toFormat(decimalPlaces, format);
    }
}
//# sourceMappingURL=Fraction.js.map