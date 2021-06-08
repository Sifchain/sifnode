"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.useField = void 0;
const reactivity_1 = require("@vue/reactivity");
const entities_1 = require("../../../entities");
const utils_1 = require("../../../utils");
function useField(amount, symbol) {
    const asset = reactivity_1.computed(() => {
        if (!symbol.value)
            return null;
        return entities_1.Asset(symbol.value);
    });
    const fieldAmount = reactivity_1.computed(() => {
        if (!asset.value || !amount.value)
            return null;
        return entities_1.AssetAmount(asset.value, utils_1.toBaseUnits(amount.value, asset.value));
    });
    return {
        fieldAmount,
        asset,
    };
}
exports.useField = useField;
//# sourceMappingURL=useField.js.map