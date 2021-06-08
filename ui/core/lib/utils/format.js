"use strict";
var __rest = (this && this.__rest) || function (s, e) {
    var t = {};
    for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p) && e.indexOf(p) < 0)
        t[p] = s[p];
    if (s != null && typeof Object.getOwnPropertySymbols === "function")
        for (var i = 0, p = Object.getOwnPropertySymbols(s); i < p.length; i++) {
            if (e.indexOf(p[i]) < 0 && Object.prototype.propertyIsEnumerable.call(s, p[i]))
                t[p[i]] = s[p[i]];
        }
    return t;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.trimMantissa = exports.format = exports.round = exports.getMantissaFromDynamicMantissa = void 0;
const Amount_1 = require("../entities/Amount");
const numbro_1 = __importDefault(require("numbro"));
const decimalShift_1 = require("./decimalShift");
function isAsset(val) {
    return !!val && typeof (val === null || val === void 0 ? void 0 : val.symbol) === "string";
}
/**
 * Takes an amount and a dynamic mantissa hash and returns the mantisaa value to use
 * @param amount amount given to format function
 * @param hash dynamic value hash to calculate mantissa from
 * @returns number of mantissa to send to formatter
 */
function getMantissaFromDynamicMantissa(amount, hash) {
    const { infinity } = hash, numHash = __rest(hash, ["infinity"]);
    const entries = Object.entries(numHash);
    entries.sort(([a], [b]) => {
        if (a > b)
            return 1;
        return -1;
    });
    for (const entry of entries) {
        const [range, mantissa] = entry;
        if (amount.lessThan(range)) {
            return mantissa;
        }
    }
    if (amount.lessThan("10000")) {
        return 2;
    }
    return infinity;
}
exports.getMantissaFromDynamicMantissa = getMantissaFromDynamicMantissa;
function round(decimal, places) {
    return decimalShift_1.decimalShift(Amount_1.Amount(decimal)
        .multiply(Amount_1.Amount(decimalShift_1.decimalShift("1", places)))
        .toBigInt() // apply rounding
        .toString(), -1 * places);
}
exports.round = round;
function isDynamicMantissa(value) {
    return typeof value !== "number";
}
function isOptionsWithFixedMantissa(options) {
    return options.shorthand || !isDynamicMantissa(options.mantissa);
}
/**
 * Options come with a dynamic or fixed mantissa. This function converts a dynamic mantissa value if it exists to a fixed number
 * @param amount
 * @param options
 * @returns
 */
function convertDynamicMantissaToFixedMantissa(amount, options) {
    if (!isOptionsWithFixedMantissa(options) &&
        typeof options.mantissa === "object") {
        return Object.assign(Object.assign({}, options), { mantissa: getMantissaFromDynamicMantissa(amount, options.mantissa) });
    }
    return options;
}
function format(_amount, _asset, _options) {
    var _a, _b;
    const amount = _amount;
    const _optionsWithDynamicMantissa = (isAsset(_asset) ? _options : _asset) || {};
    const asset = isAsset(_asset) ? _asset : undefined;
    const options = convertDynamicMantissaToFixedMantissa(amount, _optionsWithDynamicMantissa);
    // This should not happen in typed parts of the codebase
    if (typeof amount === "string") {
        // We need this in order to push developers to use the amount API right to the point at which we format values for display
        // Currently not using JSX means types are not necessarily propagated to every view so types guards
        // and there was a happy coincidence that format happened to work with a string and no asset
        //
        // We need to avoid this for the following reasons:
        //   * It encourages the status quo of not using JSX which has many poor knockon effects
        //   * One way api leads to simpler and easier to understand code
        //   * It reduces refactorability
        //   * It adds complexity to the codebase as it enables accidental amount -> string -> amount flows
        //   * It makes it more likely that developers accidentally try to format AssetAmounts as Amounts which
        //     is something this function attempts to solve using Types
        //   * It adds difficult to track down errors as strings of unknown format are passed to the format function
        //
        // Once JSX is used throughout the codebase it might be time to revisit this
        throw new Error("Amount can only take an IAmount and must NOT be a string. If you have a string and need to format it you should first convert it to an IAmount. Eg. format(Amount('100'), myformat)");
    }
    if (!amount) {
        // In theory this should not happen if we are using typescript correctly
        // This might happen due to a service response not being runtime checked
        // or in Vue because we are not using JSX templates
        console.error(`Amount "${amount}" supplied to format function is falsey`);
        return ""; // return empty string if there is an error
    }
    let decimal = asset
        ? decimalShift_1.decimalShift(amount.toBigInt().toString(), -1 * asset.decimals)
        : amount.toString();
    let postfix = (_a = options.prefix) !== null && _a !== void 0 ? _a : "";
    let prefix = (_b = options.postfix) !== null && _b !== void 0 ? _b : "";
    let space = "";
    if (options.zeroFormat && amount.equalTo("0")) {
        return options.zeroFormat;
    }
    if (options.shorthand) {
        return numbro_1.default(decimal).format(createNumbroConfig(options));
    }
    if (options.space) {
        space = " ";
    }
    if (options.mode === "percent") {
        decimal = decimalShift_1.decimalShift(decimal, 2);
        postfix = "%";
    }
    if (typeof options.mantissa === "number") {
        decimal = applyMantissa(decimal, options.mantissa);
    }
    if (options.trimMantissa) {
        decimal = trimMantissa(decimal, options.trimMantissa === "integer");
    }
    if (options.separator) {
        decimal = applySeparator(decimal);
    }
    return `${prefix}${decimal}${space}${postfix}`;
}
exports.format = format;
function trimMantissa(decimal, integer = false) {
    return decimal.replace(/(0+)$/, "").replace(/\.$/, integer ? "" : ".0");
}
exports.trimMantissa = trimMantissa;
function applySeparator(decimal) {
    const [char, mant] = decimal.split(".");
    return [char.replace(/\B(?<!\.\d*)(?=(\d{3})+(?!\d))/g, ","), mant].join(".");
}
function applyMantissa(decimal, mantissa) {
    return round(decimal, mantissa);
}
function isShorthandWithTotalLength(val) {
    return (val === null || val === void 0 ? void 0 : val.shorthand) && (val === null || val === void 0 ? void 0 : val.totalLength);
}
function createNumbroConfig(options) {
    var _a, _b, _c, _d, _e, _f, _g, _h, _j;
    return Object.assign({ forceSign: (_a = options.forceSign) !== null && _a !== void 0 ? _a : false, output: (_b = options.mode) !== null && _b !== void 0 ? _b : "number", thousandSeparated: (_c = options.separator) !== null && _c !== void 0 ? _c : false, spaceSeparated: (_d = options.space) !== null && _d !== void 0 ? _d : false, prefix: (_e = options.prefix) !== null && _e !== void 0 ? _e : "", postfix: (_f = options.postfix) !== null && _f !== void 0 ? _f : "" }, (isShorthandWithTotalLength(options)
        ? {
            average: (_g = options.shorthand) !== null && _g !== void 0 ? _g : false,
            totalLength: options.totalLength,
        }
        : {
            average: (_h = options.shorthand) !== null && _h !== void 0 ? _h : false,
            mantissa: (_j = options.mantissa) !== null && _j !== void 0 ? _j : 0,
            trimMantissa: !!options.trimMantissa,
        }));
}
//# sourceMappingURL=format.js.map