"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Amount = void 0;
const jsbi_1 = __importDefault(require("jsbi"));
const Fraction_1 = require("./fraction/Fraction");
const big_js_1 = __importDefault(require("big.js"));
const AssetAmount_1 = require("./AssetAmount");
const decimalShift_1 = require("../utils/decimalShift");
const INTEGER_REG_EX = /^[+-]?\d+$/;
const NUMBER_WITH_DECIMAL_POINT_REG_EX = /^[+-]?(\d+)?\.\d+$/;
function Amount(source) {
    // Am I a decimal number string with a period?
    if (typeof source === "string" &&
        source.match(NUMBER_WITH_DECIMAL_POINT_REG_EX)) {
        return getAmountFromDecimal(source);
    }
    // Ok so I must be an integer or something is wrong
    if (typeof source === "string" && !source.match(INTEGER_REG_EX)) {
        throw new Error(`Amount input error! string "${source}" is not numeric`);
    }
    // Our types dictate you cannot have falsey source but sometimes we
    // have casted poorly or have not validated or sanitized input
    if (!source) {
        throw new Error(`Amount input cannot be falsey given <${source}>`);
    }
    if (!(source instanceof jsbi_1.default) &&
        typeof source !== "bigint" &&
        typeof source !== "string") {
        if (AssetAmount_1.isAssetAmount(source)) {
            return source.amount;
        }
        return source;
    }
    let fraction = new Fraction_1.Fraction(source);
    const instance = {
        // We only loose precision and round when we move to BigInt for display
        toBigInt() {
            return getQuotientWithBankersRounding(fraction);
        },
        toString(detailed = true) {
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
            const string = toFraction(big.sqrt().times("100000000000000000000000").toFixed(0));
            return Amount(string).divide("100000000000000000000000");
        },
        // Internal methods need to be exposed here
        // so they can be used by another Amount in
        // toFraction and toAmount
        _fromInternal(_fraction) {
            fraction = _fraction;
            return instance;
        },
        _toInternal() {
            return fraction;
        },
    };
    return instance;
}
exports.Amount = Amount;
// quotient needs to use bankers rounding so we follow this example for bankers rounding in BigInt and apply to JSBI
//https://stackoverflow.com/questions/53752370/ecmascript-bigint-round-to-even
function getQuotientWithBankersRounding(fraction) {
    const a = fraction.numerator;
    const b = fraction.denominator;
    const aAbs = jsbi_1.default.greaterThan(a, jsbi_1.default.BigInt("0"))
        ? a
        : jsbi_1.default.multiply(jsbi_1.default.BigInt("-1"), a);
    const bAbs = jsbi_1.default.greaterThan(b, jsbi_1.default.BigInt("0"))
        ? b
        : jsbi_1.default.multiply(jsbi_1.default.BigInt("-1"), b);
    let result = jsbi_1.default.divide(aAbs, bAbs);
    const rem = jsbi_1.default.remainder(aAbs, bAbs);
    if (jsbi_1.default.greaterThan(jsbi_1.default.multiply(rem, jsbi_1.default.BigInt("2")), bAbs)) {
        result = jsbi_1.default.add(result, jsbi_1.default.BigInt("1"));
    }
    else if (jsbi_1.default.equal(jsbi_1.default.multiply(rem, jsbi_1.default.BigInt("2")), bAbs)) {
        if (jsbi_1.default.equal(jsbi_1.default.remainder(result, jsbi_1.default.BigInt("2")), jsbi_1.default.BigInt("1"))) {
            result = jsbi_1.default.add(result, jsbi_1.default.BigInt("1"));
        }
    }
    if (jsbi_1.default.greaterThan(a, jsbi_1.default.BigInt("0")) !==
        jsbi_1.default.greaterThan(b, jsbi_1.default.BigInt("0"))) {
        return jsbi_1.default.multiply(jsbi_1.default.BigInt("-1"), result);
    }
    else {
        return result;
    }
}
function getAmountFromDecimal(decimal) {
    return Amount(decimalShift_1.floorDecimal(decimalShift_1.decimalShift(decimal, 18))).divide("1000000000000000000");
}
function toFraction(a) {
    if (typeof a === "string") {
        return a.indexOf(".") < 0 ? a : Amount(a)._toInternal();
    }
    return a._toInternal();
}
// Internal helper convert to Big.js for calculating sqrts
// NOTE this looses precision to 1e24
function toBig(fraction) {
    return big_js_1.default(fraction.toFixed(24));
}
// Helper for converting a fraction to an amount.
// This uses a private API and should not be exposed
// outside of Amount
function toAmount(a) {
    return Amount("0")._fromInternal(a);
}
//# sourceMappingURL=Amount.js.map