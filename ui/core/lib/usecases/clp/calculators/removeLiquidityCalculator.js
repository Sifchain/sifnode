"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.useRemoveLiquidityCalculator = void 0;
const entities_1 = require("../../../entities");
const formulae_1 = require("../../../entities/formulae");
const format_1 = require("../../../utils/format");
const addLiquidityCalculator_1 = require("./addLiquidityCalculator");
function useRemoveLiquidityCalculator(input) {
    // this function needs to be refactored so
    const externalAsset = (() => {
        if (!input.externalAssetSymbol.value)
            return null;
        return entities_1.Asset(input.externalAssetSymbol.value);
    })();
    const nativeAsset = (() => {
        if (!input.nativeAssetSymbol.value)
            return null;
        return entities_1.Asset(input.nativeAssetSymbol.value);
    })();
    const liquidityPool = (() => {
        var _a;
        if (!nativeAsset || !externalAsset)
            return null;
        // Find pool from poolFinder
        const pool = input.poolFinder(externalAsset, nativeAsset);
        return (_a = pool === null || pool === void 0 ? void 0 : pool.value) !== null && _a !== void 0 ? _a : null;
    })();
    const poolUnits = (() => {
        if (!liquidityPool)
            return null;
        return liquidityPool.poolUnits;
    })();
    const wBasisPoints = (() => {
        if (!input.wBasisPoints.value)
            return null;
        return entities_1.Amount(input.wBasisPoints.value);
    })();
    const asymmetry = (() => {
        if (!input.asymmetry.value)
            return null;
        return entities_1.Amount(input.asymmetry.value);
    })();
    const nativeAssetBalance = (() => {
        var _a;
        if (!liquidityPool)
            return null;
        return ((_a = liquidityPool.amounts.find((a) => a.symbol === input.nativeAssetSymbol.value)) !== null && _a !== void 0 ? _a : null);
    })();
    const externalAssetBalance = (() => {
        var _a;
        if (!liquidityPool)
            return null;
        return ((_a = liquidityPool.amounts.find((a) => a.symbol === input.externalAssetSymbol.value)) !== null && _a !== void 0 ? _a : null);
    })();
    const lpUnits = (() => {
        if (!input.liquidityProvider.value)
            return null;
        return input.liquidityProvider.value.units;
    })();
    const hasLiquidity = (() => {
        if (!lpUnits)
            return false;
        return lpUnits.greaterThan("0");
    })();
    const withdrawalAmounts = (() => {
        if (!poolUnits ||
            !nativeAssetBalance ||
            !externalAssetBalance ||
            !lpUnits ||
            !wBasisPoints ||
            !asymmetry ||
            !externalAsset ||
            !nativeAsset)
            return null;
        const { withdrawExternalAssetAmount, withdrawNativeAssetAmount, } = formulae_1.calculateWithdrawal({
            poolUnits,
            nativeAssetBalance,
            externalAssetBalance,
            lpUnits,
            wBasisPoints,
            asymmetry: asymmetry,
        });
        return {
            hasLiquidity,
            withdrawExternalAssetAmount: entities_1.AssetAmount(externalAsset, withdrawExternalAssetAmount),
            withdrawNativeAssetAmount: entities_1.AssetAmount(nativeAsset, withdrawNativeAssetAmount),
        };
    })();
    const state = (() => {
        if (!input.externalAssetSymbol.value || !input.nativeAssetSymbol.value)
            return addLiquidityCalculator_1.PoolState.SELECT_TOKENS;
        if (!(wBasisPoints === null || wBasisPoints === void 0 ? void 0 : wBasisPoints.greaterThan("0")))
            return addLiquidityCalculator_1.PoolState.ZERO_AMOUNTS;
        if (!hasLiquidity)
            return addLiquidityCalculator_1.PoolState.NO_LIQUIDITY;
        if (!lpUnits) {
            return addLiquidityCalculator_1.PoolState.INSUFFICIENT_FUNDS;
        }
        return addLiquidityCalculator_1.PoolState.VALID_INPUT;
    })();
    const withdrawExternalAssetAmountMessage = (() => {
        if (!withdrawalAmounts)
            return "";
        const assetAmount = withdrawalAmounts === null || withdrawalAmounts === void 0 ? void 0 : withdrawalAmounts.withdrawExternalAssetAmount;
        return format_1.format(assetAmount.amount, assetAmount.asset, {
            mantissa: 6,
        });
    })();
    const withdrawNativeAssetAmountMessage = (() => {
        if (!withdrawalAmounts)
            return "";
        const assetAmount = withdrawalAmounts === null || withdrawalAmounts === void 0 ? void 0 : withdrawalAmounts.withdrawNativeAssetAmount;
        return format_1.format(assetAmount.amount, assetAmount.asset, {
            mantissa: 6,
        });
    })();
    return {
        withdrawExternalAssetAmount: withdrawExternalAssetAmountMessage,
        withdrawNativeAssetAmount: withdrawNativeAssetAmountMessage,
        state,
    };
}
exports.useRemoveLiquidityCalculator = useRemoveLiquidityCalculator;
//# sourceMappingURL=removeLiquidityCalculator.js.map