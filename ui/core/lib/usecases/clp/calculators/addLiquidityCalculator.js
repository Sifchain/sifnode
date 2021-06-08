"use strict";
// TODO remove refs dependency and move to `actions/clp/calculateAddLiquidity`
Object.defineProperty(exports, "__esModule", { value: true });
exports.usePoolCalculator = exports.PoolState = void 0;
const reactivity_1 = require("@vue/reactivity");
const entities_1 = require("../../../entities");
const entities_2 = require("../../../entities");
const format_1 = require("../../../utils/format");
const useField_1 = require("./useField");
const utils_1 = require("./utils");
var PoolState;
(function (PoolState) {
    PoolState[PoolState["SELECT_TOKENS"] = 0] = "SELECT_TOKENS";
    PoolState[PoolState["ZERO_AMOUNTS"] = 1] = "ZERO_AMOUNTS";
    PoolState[PoolState["INSUFFICIENT_FUNDS"] = 2] = "INSUFFICIENT_FUNDS";
    PoolState[PoolState["VALID_INPUT"] = 3] = "VALID_INPUT";
    PoolState[PoolState["NO_LIQUIDITY"] = 4] = "NO_LIQUIDITY";
    PoolState[PoolState["ZERO_AMOUNTS_NEW_POOL"] = 5] = "ZERO_AMOUNTS_NEW_POOL";
})(PoolState = exports.PoolState || (exports.PoolState = {}));
function usePoolCalculator(input) {
    const tokenAField = useField_1.useField(input.tokenAAmount, input.tokenASymbol);
    const tokenBField = useField_1.useField(input.tokenBAmount, input.tokenBSymbol);
    const balanceMap = utils_1.useBalances(input.balances);
    const preExistingPool = reactivity_1.computed(() => {
        if (!tokenAField.asset.value || !tokenBField.asset.value) {
            return null;
        }
        // Find pool from poolFinder
        const pool = input.poolFinder(tokenAField.asset.value.symbol, tokenBField.asset.value.symbol);
        return (pool === null || pool === void 0 ? void 0 : pool.value) || null;
    });
    const assetA = reactivity_1.computed(() => {
        if (!input.tokenASymbol.value) {
            return null;
        }
        return entities_1.Asset.get(input.tokenASymbol.value);
    });
    const assetB = reactivity_1.computed(() => {
        if (!input.tokenBSymbol.value) {
            return null;
        }
        return entities_1.Asset.get(input.tokenBSymbol.value);
    });
    const tokenABalance = reactivity_1.computed(() => {
        var _a, _b;
        if (!tokenAField.fieldAmount.value || !tokenAField.asset.value) {
            return null;
        }
        if (preExistingPool.value) {
            return input.tokenASymbol.value
                ? (_a = balanceMap.value.get(input.tokenASymbol.value)) !== null && _a !== void 0 ? _a : entities_1.AssetAmount(tokenAField.asset.value, "0") : null;
        }
        else {
            return input.tokenASymbol.value
                ? (_b = balanceMap.value.get(input.tokenASymbol.value)) !== null && _b !== void 0 ? _b : null : null;
        }
    });
    const tokenBBalance = reactivity_1.computed(() => {
        var _a;
        return input.tokenBSymbol.value
            ? (_a = balanceMap.value.get(input.tokenBSymbol.value)) !== null && _a !== void 0 ? _a : null : null;
    });
    const fromBalanceOverdrawn = reactivity_1.computed(() => {
        var _a;
        return !((_a = tokenABalance.value) === null || _a === void 0 ? void 0 : _a.greaterThanOrEqual(tokenAField.fieldAmount.value || "0"));
    });
    const toBalanceOverdrawn = reactivity_1.computed(() => {
        var _a;
        return !((_a = tokenBBalance.value) === null || _a === void 0 ? void 0 : _a.greaterThanOrEqual(tokenBField.fieldAmount.value || "0"));
    });
    const liquidityPool = reactivity_1.computed(() => {
        if (preExistingPool.value) {
            return preExistingPool.value;
        }
        if (!tokenAField.fieldAmount.value ||
            !tokenBField.fieldAmount.value ||
            !tokenAField.asset.value ||
            !tokenBField.asset.value) {
            return null;
        }
        return entities_1.Pool(entities_1.AssetAmount(tokenAField.asset.value, "0"), entities_1.AssetAmount(tokenBField.asset.value, "0"));
    });
    // pool units for this prospective transaction [total, newUnits]
    const provisionedPoolUnitsArray = reactivity_1.computed(() => {
        if (!liquidityPool.value ||
            !tokenBField.fieldAmount.value ||
            !tokenAField.fieldAmount.value) {
            return [entities_2.Amount("0"), entities_2.Amount("0")];
        }
        return liquidityPool.value.calculatePoolUnits(tokenBField.fieldAmount.value, tokenAField.fieldAmount.value);
    });
    // pool units from the perspective of the liquidity provider
    const liquidityProviderPoolUnitsArray = reactivity_1.computed(() => {
        if (!provisionedPoolUnitsArray.value)
            return [entities_2.Amount("0"), entities_2.Amount("0")];
        const [totalPoolUnits, newUnits] = provisionedPoolUnitsArray.value;
        // if this user already has pool units include those in the newUnits
        const totalLiquidityProviderUnits = input.liquidityProvider.value
            ? input.liquidityProvider.value.units.add(newUnits)
            : newUnits;
        return [totalPoolUnits, totalLiquidityProviderUnits];
    });
    const totalPoolUnits = reactivity_1.computed(() => liquidityProviderPoolUnitsArray.value[0].toBigInt().toString());
    const totalLiquidityProviderUnits = reactivity_1.computed(() => liquidityProviderPoolUnitsArray.value[1].toBigInt().toString());
    const shareOfPool = reactivity_1.computed(() => {
        if (!liquidityProviderPoolUnitsArray.value)
            return entities_2.Amount("0");
        const [units, lpUnits] = liquidityProviderPoolUnitsArray.value;
        // shareOfPool should be 0 if units and lpUnits are zero
        if (units.equalTo("0") && lpUnits.equalTo("0"))
            return entities_2.Amount("0");
        // if no units lp owns 100% of pool
        return units.equalTo("0") ? entities_2.Amount("1") : lpUnits.divide(units);
    });
    const shareOfPoolPercent = reactivity_1.computed(() => {
        if (shareOfPool.value.multiply("10000").lessThan("1"))
            return "< 0.01%";
        return `${format_1.format(shareOfPool.value, {
            mantissa: 2,
            mode: "percent",
        })}`;
    });
    const poolAmounts = reactivity_1.computed(() => {
        if (!preExistingPool.value || !tokenAField.asset.value) {
            return null;
        }
        if (!preExistingPool.value.contains(tokenAField.asset.value))
            return null;
        const externalBalance = preExistingPool.value.getAmount(tokenAField.asset.value);
        const nativeBalance = preExistingPool.value.getAmount("rowan");
        return [nativeBalance, externalBalance];
    });
    // external_balance / native_balance
    const aPerBRatio = reactivity_1.computed(() => {
        if (!poolAmounts.value)
            return 0;
        const [native, external] = poolAmounts.value;
        const derivedNative = native.toDerived();
        const derivedExternal = external.toDerived();
        return derivedExternal.divide(derivedNative);
    });
    const aPerBRatioMessage = reactivity_1.computed(() => {
        if (!aPerBRatio.value) {
            return "N/A";
        }
        return format_1.format(aPerBRatio.value, { mantissa: 8 });
    });
    // native_balance / external_balance
    const bPerARatio = reactivity_1.computed(() => {
        if (!poolAmounts.value)
            return 0;
        const [native, external] = poolAmounts.value;
        const derivedNative = native.toDerived();
        const derivedExternal = external.toDerived();
        return derivedNative.divide(derivedExternal);
    });
    const bPerARatioMessage = reactivity_1.computed(() => {
        if (!bPerARatio.value) {
            return "N/A";
        }
        return format_1.format(bPerARatio.value, { mantissa: 8 });
    });
    // Price Impact and Pool Share:
    // (external_balance + external_added) / (native_balance + native_added)
    const aPerBRatioProjected = reactivity_1.computed(() => {
        if (!poolAmounts.value ||
            !tokenAField.fieldAmount.value ||
            !tokenBField.fieldAmount.value)
            return null;
        const [native, external] = poolAmounts.value;
        const derivedNative = native.toDerived();
        const derivedExternal = external.toDerived();
        const externalAdded = tokenAField.fieldAmount.value.toDerived();
        const nativeAdded = tokenBField.fieldAmount.value.toDerived();
        return derivedExternal
            .add(externalAdded)
            .divide(derivedNative.add(nativeAdded));
    });
    const aPerBRatioProjectedMessage = reactivity_1.computed(() => {
        if (!aPerBRatioProjected.value) {
            return "N/A";
        }
        return format_1.format(aPerBRatioProjected.value, { mantissa: 8 });
    });
    // Price Impact and Pool Share:
    // (native_balance + native_added)/(external_balance + external_added)
    const bPerARatioProjected = reactivity_1.computed(() => {
        if (!poolAmounts.value ||
            !tokenAField.fieldAmount.value ||
            !tokenBField.fieldAmount.value)
            return null;
        const [native, external] = poolAmounts.value;
        const derivedNative = native.toDerived();
        const derivedExternal = external.toDerived();
        const externalAdded = tokenAField.fieldAmount.value.toDerived();
        const nativeAdded = tokenBField.fieldAmount.value.toDerived();
        return derivedNative
            .add(nativeAdded)
            .divide(derivedExternal.add(externalAdded));
    });
    const bPerARatioProjectedMessage = reactivity_1.computed(() => {
        if (!bPerARatioProjected.value) {
            return "N/A";
        }
        return format_1.format(bPerARatioProjected.value, { mantissa: 8 });
    });
    reactivity_1.effect(() => {
        var _a, _b;
        // if in guided mode
        // calculate the price ratio of A / B
        // Only activates when it is a preexisting pool
        if (assetA.value &&
            assetB.value &&
            input.asyncPooling.value &&
            preExistingPool.value &&
            input.lastFocusedTokenField.value !== null) {
            if (bPerARatio === null ||
                aPerBRatio === null ||
                !assetA.value ||
                !assetB.value) {
                return null;
            }
            const assetAmountA = entities_1.AssetAmount(assetA.value, ((_a = tokenAField.fieldAmount) === null || _a === void 0 ? void 0 : _a.value) || "0");
            const assetAmountB = entities_1.AssetAmount(assetB.value, ((_b = tokenBField.fieldAmount) === null || _b === void 0 ? void 0 : _b.value) || "0");
            if (input.lastFocusedTokenField.value === "A") {
                input.tokenBAmount.value = format_1.format(assetAmountA.toDerived().multiply(bPerARatio.value || "0"), { mantissa: 5 });
            }
            else if (input.lastFocusedTokenField.value === "B") {
                input.tokenAAmount.value = format_1.format(assetAmountB.toDerived().multiply(aPerBRatio.value || "0"), { mantissa: 5 });
            }
        }
    });
    const state = reactivity_1.computed(() => {
        // Select Tokens
        const aSymbolNotSelected = !input.tokenASymbol.value;
        const bSymbolNotSelected = !input.tokenBSymbol.value;
        if (aSymbolNotSelected || bSymbolNotSelected) {
            return PoolState.SELECT_TOKENS;
        }
        // Zero amounts
        const aAmount = tokenAField.fieldAmount.value;
        const bAmount = tokenBField.fieldAmount.value;
        const aAmountIsZeroOrFalsy = !aAmount || aAmount.equalTo("0");
        const bAmountIsZeroOrFalsy = !bAmount || bAmount.equalTo("0");
        if (!preExistingPool.value &&
            (aAmountIsZeroOrFalsy || bAmountIsZeroOrFalsy)) {
            return PoolState.ZERO_AMOUNTS_NEW_POOL;
        }
        if (aAmountIsZeroOrFalsy && bAmountIsZeroOrFalsy) {
            return PoolState.ZERO_AMOUNTS;
        }
        // Insufficient Funds
        if (fromBalanceOverdrawn.value || toBalanceOverdrawn.value) {
            return PoolState.INSUFFICIENT_FUNDS;
        }
        // Valid yay!
        return PoolState.VALID_INPUT;
    });
    return {
        state,
        aPerBRatioMessage,
        bPerARatioMessage,
        aPerBRatioProjectedMessage,
        bPerARatioProjectedMessage,
        shareOfPool,
        shareOfPoolPercent,
        preExistingPool,
        totalLiquidityProviderUnits,
        totalPoolUnits,
        poolAmounts,
        tokenAFieldAmount: tokenAField.fieldAmount,
        tokenBFieldAmount: tokenBField.fieldAmount,
    };
}
exports.usePoolCalculator = usePoolCalculator;
//# sourceMappingURL=addLiquidityCalculator.js.map