"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Pair = void 0;
function Pair(a, b) {
    const amounts = [a, b];
    return {
        amounts,
        otherAsset(asset) {
            const otherAsset = amounts.find((amount) => amount.symbol !== asset.symbol);
            if (!otherAsset)
                throw new Error("Asset doesnt exist in pair");
            return otherAsset;
        },
        symbol() {
            return amounts
                .map((a) => a.symbol)
                .sort()
                .join("_");
        },
        contains(...assets) {
            const local = amounts.map((a) => a.symbol);
            const other = assets.map((a) => a.symbol);
            return !!local.find((s) => other.includes(s));
        },
        getAmount(asset) {
            const assetSymbol = typeof asset === "string" ? asset : asset.symbol;
            const found = this.amounts.find((amount) => {
                return amount.symbol === assetSymbol;
            });
            if (!found)
                throw new Error(`Asset ${assetSymbol} doesnt exist in pair`);
            return found;
        },
        toString() {
            return amounts.map((a) => a.toString()).join(" | ");
        },
    };
}
exports.Pair = Pair;
//# sourceMappingURL=Pair.js.map