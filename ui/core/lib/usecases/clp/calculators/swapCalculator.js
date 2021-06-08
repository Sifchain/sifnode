"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.useSwapCalculator = exports.SwapState = void 0;
const reactivity_1 = require("@vue/reactivity");
const entities_1 = require("../../../entities");
const useField_1 = require("./useField");
const utils_1 = require("./utils");
const format_1 = require("../../../utils/format");
var SwapState;
(function (SwapState) {
    SwapState[SwapState["SELECT_TOKENS"] = 0] = "SELECT_TOKENS";
    SwapState[SwapState["ZERO_AMOUNTS"] = 1] = "ZERO_AMOUNTS";
    SwapState[SwapState["INSUFFICIENT_FUNDS"] = 2] = "INSUFFICIENT_FUNDS";
    SwapState[SwapState["VALID_INPUT"] = 3] = "VALID_INPUT";
    SwapState[SwapState["INVALID_AMOUNT"] = 4] = "INVALID_AMOUNT";
    SwapState[SwapState["INSUFFICIENT_LIQUIDITY"] = 5] = "INSUFFICIENT_LIQUIDITY";
})(SwapState = exports.SwapState || (exports.SwapState = {}));
function calculateFormattedPriceImpact(pair, amount) {
    return format_1.format(pair.calcPriceImpact(amount), {
        mantissa: 6,
        trimMantissa: true,
    });
}
function calculateFormattedProviderFee(pair, amount) {
    return format_1.format(pair.calcProviderFee(amount).amount, pair.calcProviderFee(amount).asset, { mantissa: 5, trimMantissa: true });
}
// TODO: make swap calculator only generate Fractions/Amounts that get stringified in the view
function useSwapCalculator(input) {
    // extracting selectedField so we can use it without tracking its change
    let selectedField = null;
    reactivity_1.effect(() => (selectedField = input.selectedField.value));
    // We use a composite pool pair to work out rates
    const pool = reactivity_1.computed(() => {
        var _a, _b, _c, _d;
        if (!input.fromSymbol.value || !input.toSymbol.value)
            return null;
        if (input.fromSymbol.value === "rowan") {
            return (_b = (_a = input.poolFinder(input.toSymbol.value, "rowan")) === null || _a === void 0 ? void 0 : _a.value) !== null && _b !== void 0 ? _b : null;
        }
        if (input.toSymbol.value === "rowan") {
            return (_d = (_c = input.poolFinder(input.fromSymbol.value, "rowan")) === null || _c === void 0 ? void 0 : _c.value) !== null && _d !== void 0 ? _d : null;
        }
        const fromPair = input.poolFinder(input.fromSymbol.value, "rowan");
        const toPair = input.poolFinder(input.toSymbol.value, "rowan");
        if (!fromPair || !toPair)
            return null;
        return entities_1.CompositePool(fromPair.value, toPair.value);
    });
    // Get the balance of the from the users account
    const balance = reactivity_1.computed(() => {
        var _a;
        const balanceMap = utils_1.useBalances(input.balances);
        return input.fromSymbol.value
            ? (_a = balanceMap.value.get(input.fromSymbol.value)) !== null && _a !== void 0 ? _a : null : null;
    });
    // Get field amounts as domain objects
    const fromField = useField_1.useField(input.fromAmount, input.fromSymbol);
    const toField = useField_1.useField(input.toAmount, input.toSymbol);
    // Create a price message eg. 10.123 ATK per BTK
    const priceMessage = reactivity_1.computed(() => {
        if (!fromField.fieldAmount.value ||
            fromField.fieldAmount.value.equalTo("0") ||
            !pool.value) {
            return "";
        }
        const amount = fromField.fieldAmount.value;
        const pair = pool.value;
        const swapResult = pair.calcSwapResult(amount);
        // to get ratio needs to be divided by amount as input by user
        const amountAsInput = format_1.format(amount.amount, amount.asset);
        return `${format_1.format(swapResult.divide(amountAsInput), swapResult.asset, {
            mantissa: 6,
        })} ${swapResult.label} per ${amount.label}`;
    });
    // Selected field changes when the user changes the field selection
    // If the selected field is the "tokenA" field and something changes we change the "tokenB" input value
    // If the selected field is the "tokenB" field and something changes we change the "tokenA" input value
    // Changing the "from" field recalculates the "to" amount
    const swapResult = reactivity_1.ref(null);
    reactivity_1.effect(() => {
        if (pool.value &&
            fromField.asset.value &&
            fromField.fieldAmount.value &&
            pool.value.contains(fromField.asset.value) &&
            selectedField === "from") {
            swapResult.value = pool.value.calcSwapResult(fromField.fieldAmount.value);
            const toAmountValue = format_1.format(swapResult.value.amount, swapResult.value.asset, {
                mantissa: 10,
                trimMantissa: true,
            });
            input.toAmount.value = toAmountValue;
        }
    });
    // Changing the "to" field recalculates the "from" amount
    const reverseSwapResult = reactivity_1.ref(null);
    reactivity_1.effect(() => {
        if (pool.value &&
            toField.asset.value &&
            toField.fieldAmount.value &&
            pool.value.contains(toField.asset.value) &&
            selectedField === "to") {
            reverseSwapResult.value = pool.value.calcReverseSwapResult(toField.fieldAmount.value);
            // Internally trigger calulations based off swapResult as this is how we
            // work out priceImpact, providerFee, minimumReceived
            swapResult.value = pool.value.calcSwapResult(reverseSwapResult.value);
            input.fromAmount.value = utils_1.trimZeros(format_1.format(reverseSwapResult.value.amount, reverseSwapResult.value.asset, {
                mantissa: 8,
            }));
        }
    });
    // Format input amount on blur
    reactivity_1.effect(() => {
        if (input.selectedField.value === null && input.toAmount.value) {
            input.toAmount.value = utils_1.trimZeros(input.toAmount.value);
        }
    });
    // Format input amount on blur
    reactivity_1.effect(() => {
        if (input.selectedField.value === null && input.fromAmount.value) {
            input.fromAmount.value = utils_1.trimZeros(input.fromAmount.value);
        }
    });
    // Cache pool contains asset for reuse as is a little
    const poolContainsFromAsset = reactivity_1.computed(() => {
        if (!fromField.asset.value || !pool.value)
            return false;
        return pool.value.contains(fromField.asset.value);
    });
    const priceImpact = reactivity_1.computed(() => {
        if (!pool.value ||
            !fromField.asset.value ||
            !fromField.fieldAmount.value ||
            !poolContainsFromAsset.value)
            return null;
        return calculateFormattedPriceImpact(pool.value, fromField.fieldAmount.value);
    });
    const providerFee = reactivity_1.computed(() => {
        if (!pool.value ||
            !fromField.asset.value ||
            !fromField.fieldAmount.value ||
            !poolContainsFromAsset.value)
            return null;
        return calculateFormattedProviderFee(pool.value, fromField.fieldAmount.value);
    });
    // minimumReceived
    const minimumReceived = reactivity_1.computed(() => {
        if (!input.slippage.value || !toField.asset.value || !swapResult.value)
            return null;
        const slippage = entities_1.Amount(input.slippage.value);
        const minAmount = entities_1.Amount("1")
            .subtract(slippage.divide(entities_1.Amount("100")))
            .multiply(swapResult.value);
        return entities_1.AssetAmount(toField.asset.value, minAmount);
    });
    // Derive state
    const state = reactivity_1.computed(() => {
        var _a, _b, _c;
        // SwapState.INSUFFICIENT_LIQUIDITY is probably better here
        if (!pool.value)
            return SwapState.SELECT_TOKENS;
        const fromTokenLiquidity = pool.value.amounts.find((amount) => { var _a; return amount.asset.symbol === ((_a = fromField.asset.value) === null || _a === void 0 ? void 0 : _a.symbol); });
        const toTokenLiquidity = pool.value.amounts.find((amount) => { var _a; return amount.asset.symbol === ((_a = toField.asset.value) === null || _a === void 0 ? void 0 : _a.symbol); });
        if (!fromTokenLiquidity ||
            !toTokenLiquidity ||
            !fromField.fieldAmount.value ||
            !toField.fieldAmount.value ||
            (((_a = fromField.fieldAmount.value) === null || _a === void 0 ? void 0 : _a.equalTo("0")) && ((_b = toField.fieldAmount.value) === null || _b === void 0 ? void 0 : _b.equalTo("0")))) {
            return SwapState.ZERO_AMOUNTS;
        }
        if (toField.fieldAmount.value.greaterThan("0") &&
            fromField.fieldAmount.value.equalTo("0")) {
            return SwapState.INVALID_AMOUNT;
        }
        if (!((_c = balance.value) === null || _c === void 0 ? void 0 : _c.greaterThanOrEqual(fromField.fieldAmount.value || "0")))
            return SwapState.INSUFFICIENT_FUNDS;
        if (fromTokenLiquidity.lessThan(fromField.fieldAmount.value) ||
            toTokenLiquidity.lessThan(toField.fieldAmount.value)) {
            return SwapState.INSUFFICIENT_LIQUIDITY;
        }
        return SwapState.VALID_INPUT;
    });
    return {
        priceMessage,
        state,
        fromFieldAmount: fromField.fieldAmount,
        toFieldAmount: toField.fieldAmount,
        toAmount: input.toAmount,
        fromAmount: input.fromAmount,
        priceImpact,
        providerFee,
        minimumReceived,
        swapResult,
        reverseSwapResult,
    };
}
exports.useSwapCalculator = useSwapCalculator;
//# sourceMappingURL=swapCalculator.js.map