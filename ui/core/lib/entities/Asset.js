"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Asset = void 0;
const assetMap = new Map();
function isAsset(value) {
    return (typeof (value === null || value === void 0 ? void 0 : value.symbol) === "string" && typeof (value === null || value === void 0 ? void 0 : value.decimals) === "number");
}
function Asset(assetOrSymbol) {
    // If it is an asset then cache it and return it
    if (isAsset(assetOrSymbol)) {
        assetMap.set(assetOrSymbol.symbol.toLowerCase(), assetOrSymbol);
        return assetOrSymbol;
    }
    // Return it from cache
    const found = assetOrSymbol
        ? assetMap.get(assetOrSymbol.toLowerCase())
        : false;
    if (!found) {
        throw new Error(`Attempt to retrieve the asset with key "${assetOrSymbol}" before it had been cached.`);
    }
    return found;
}
exports.Asset = Asset;
// XXX:Legacy
Asset.set = (symbol, asset) => {
    Asset(asset); // assuming symbol is same
};
// XXX:Legacy
Asset.get = (symbol) => {
    return Asset(symbol);
};
//# sourceMappingURL=Asset.js.map