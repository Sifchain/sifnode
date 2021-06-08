"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.toPool = void 0;
const entities_1 = require("../../../entities");
function getAssetOrNull(symbol) {
    try {
        return entities_1.Asset.get(symbol);
    }
    catch (err) {
        return null;
    }
}
exports.toPool = (nativeAsset) => (poolData) => {
    const externalAssetSymbol = poolData.external_asset.symbol;
    const externalAsset = getAssetOrNull(externalAssetSymbol);
    // If we are not configured to handle this external asset
    // the pool is invalid so we ignore it
    if (!externalAsset)
        return null;
    return entities_1.Pool(entities_1.AssetAmount(nativeAsset, poolData.native_asset_balance), entities_1.AssetAmount(externalAsset, poolData.external_asset_balance), entities_1.Amount(poolData.pool_units));
};
//# sourceMappingURL=toPool.js.map