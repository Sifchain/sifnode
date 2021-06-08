"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.useBalances = exports.trimZeros = exports.assetPriceMessage = void 0;
const reactivity_1 = require("@vue/reactivity");
const format_1 = require("../../../utils/format");
function assetPriceMessage(amount, pair, decimals) {
    if (!pair || !amount || amount.equalTo("0"))
        return "";
    const swapResult = pair.calcSwapResult(amount);
    return `${format_1.format(swapResult.divide(amount), {
        mantissa: decimals,
    })} ${swapResult.label} per ${amount.label}`;
}
exports.assetPriceMessage = assetPriceMessage;
function trimZeros(amount) {
    if (amount.indexOf(".") === -1)
        return `${amount}.0`;
    const tenDecimalsMax = parseFloat(amount).toFixed(10);
    return tenDecimalsMax.replace(/0+$/, "").replace(/\.$/, ".0");
}
exports.trimZeros = trimZeros;
function useBalances(balances) {
    return reactivity_1.computed(() => {
        const map = new Map();
        for (const item of balances.value) {
            map.set(item.asset.symbol, item);
        }
        return map;
    });
}
exports.useBalances = useBalances;
//# sourceMappingURL=utils.js.map