"use strict";
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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Fraction = exports.isFraction = exports.parseBigintIsh = exports._1000 = exports._997 = exports._100 = exports.TEN = exports.FIVE = exports.THREE = exports.TWO = exports.ONE = exports.ZERO = exports.Rounding = void 0;
const tiny_invariant_1 = __importDefault(require("tiny-invariant"));
const jsbi_1 = __importDefault(require("jsbi"));
const decimal_js_light_1 = __importDefault(require("decimal.js-light"));
const big_js_1 = __importDefault(require("big.js"));
const toformat_1 = __importDefault(require("toformat"));
var Rounding;
(function (Rounding) {
    Rounding[Rounding["ROUND_DOWN"] = 0] = "ROUND_DOWN";
    Rounding[Rounding["ROUND_HALF_UP"] = 1] = "ROUND_HALF_UP";
    Rounding[Rounding["ROUND_UP"] = 2] = "ROUND_UP";
})(Rounding = exports.Rounding || (exports.Rounding = {}));
exports.ZERO = jsbi_1.default.BigInt(0);
exports.ONE = jsbi_1.default.BigInt(1);
exports.TWO = jsbi_1.default.BigInt(2);
exports.THREE = jsbi_1.default.BigInt(3);
exports.FIVE = jsbi_1.default.BigInt(5);
exports.TEN = jsbi_1.default.BigInt(10);
exports._100 = jsbi_1.default.BigInt(100);
exports._997 = jsbi_1.default.BigInt(997);
exports._1000 = jsbi_1.default.BigInt(1000);
function parseBigintIsh(bigintIsh) {
    return bigintIsh instanceof jsbi_1.default
        ? bigintIsh
        : typeof bigintIsh === "bigint"
            ? jsbi_1.default.BigInt(bigintIsh.toString())
            : jsbi_1.default.BigInt(bigintIsh);
}
exports.parseBigintIsh = parseBigintIsh;
const Decimal = toformat_1.default(decimal_js_light_1.default);
const Big = toformat_1.default(big_js_1.default);
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
function isFraction(value) {
    return value.quotient instanceof jsbi_1.default;
}
exports.isFraction = isFraction;
const ensureFraction = (other) => {
    return other instanceof Fraction || isFraction(other)
        ? other
        : new Fraction(parseBigintIsh(other));
};
class Fraction {
    constructor(numerator, denominator = exports.ONE) {
        this.numerator = parseBigintIsh(numerator);
        this.denominator = parseBigintIsh(denominator);
    }
    // performs floor division
    get quotient() {
        return jsbi_1.default.divide(this.numerator, this.denominator);
    }
    // remainder after floor division
    get remainder() {
        return new Fraction(jsbi_1.default.remainder(this.numerator, this.denominator), this.denominator);
    }
    invert() {
        return new Fraction(this.denominator, this.numerator);
    }
    add(other) {
        const otherParsed = ensureFraction(other);
        if (jsbi_1.default.equal(this.denominator, otherParsed.denominator)) {
            return new Fraction(jsbi_1.default.add(this.numerator, otherParsed.numerator), this.denominator);
        }
        return new Fraction(jsbi_1.default.add(jsbi_1.default.multiply(this.numerator, otherParsed.denominator), jsbi_1.default.multiply(otherParsed.numerator, this.denominator)), jsbi_1.default.multiply(this.denominator, otherParsed.denominator));
    }
    subtract(other) {
        const otherParsed = ensureFraction(other);
        if (jsbi_1.default.equal(this.denominator, otherParsed.denominator)) {
            return new Fraction(jsbi_1.default.subtract(this.numerator, otherParsed.numerator), this.denominator);
        }
        return new Fraction(jsbi_1.default.subtract(jsbi_1.default.multiply(this.numerator, otherParsed.denominator), jsbi_1.default.multiply(otherParsed.numerator, this.denominator)), jsbi_1.default.multiply(this.denominator, otherParsed.denominator));
    }
    lessThan(other) {
        const otherParsed = ensureFraction(other);
        return jsbi_1.default.lessThan(jsbi_1.default.multiply(this.numerator, otherParsed.denominator), jsbi_1.default.multiply(otherParsed.numerator, this.denominator));
    }
    lessThanOrEqual(other) {
        return this.lessThan(other) || this.equalTo(other);
    }
    equalTo(other) {
        const otherParsed = ensureFraction(other);
        return jsbi_1.default.equal(jsbi_1.default.multiply(this.numerator, otherParsed.denominator), jsbi_1.default.multiply(otherParsed.numerator, this.denominator));
    }
    greaterThan(other) {
        const otherParsed = ensureFraction(other);
        return jsbi_1.default.greaterThan(jsbi_1.default.multiply(this.numerator, otherParsed.denominator), jsbi_1.default.multiply(otherParsed.numerator, this.denominator));
    }
    greaterThanOrEqual(other) {
        return this.greaterThan(other) || this.equalTo(other);
    }
    multiply(other) {
        const otherParsed = ensureFraction(other);
        return new Fraction(jsbi_1.default.multiply(this.numerator, otherParsed.numerator), jsbi_1.default.multiply(this.denominator, otherParsed.denominator));
    }
    divide(other) {
        const otherParsed = ensureFraction(other);
        return new Fraction(jsbi_1.default.multiply(this.numerator, otherParsed.denominator), jsbi_1.default.multiply(this.denominator, otherParsed.numerator));
    }
    toSignificant(significantDigits, format = { groupSeparator: "" }, rounding = Rounding.ROUND_HALF_UP) {
        tiny_invariant_1.default(Number.isInteger(significantDigits), `${significantDigits} is not an integer.`);
        tiny_invariant_1.default(significantDigits > 0, `${significantDigits} is not positive.`);
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
        tiny_invariant_1.default(Number.isInteger(decimalPlaces), `${decimalPlaces} is not an integer.`);
        tiny_invariant_1.default(decimalPlaces >= 0, `${decimalPlaces} is negative.`);
        Big.DP = decimalPlaces;
        Big.RM = toFixedRounding[rounding];
        return new Big(this.numerator.toString())
            .div(this.denominator.toString())
            .toFormat(decimalPlaces, format);
    }
}
exports.Fraction = Fraction;
//# sourceMappingURL=Fraction.js.map