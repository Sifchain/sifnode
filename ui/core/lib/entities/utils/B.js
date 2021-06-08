"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const jsbi_1 = __importDefault(require("jsbi"));
// Convenience method for converting a floating point number with decimals
// to a bigint representation according to a number of decimals
function B(num, dec = 18) {
    const numstr = typeof num !== "string" ? num.toFixed(dec) : num;
    const [s, m = "0", huh] = numstr.split(".");
    if (typeof huh !== "undefined")
        throw new Error("Invalid number string");
    const mm = m.length > dec ? m.slice(0, dec) : m.padEnd(dec, "0");
    const n = [s, mm].join("").replace(/^0+/, "");
    return jsbi_1.default.BigInt(n);
}
exports.default = B;
//# sourceMappingURL=B.js.map