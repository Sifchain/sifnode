import { Ref, computed, effect, ref } from "@vue/reactivity";
import {
  Asset,
  AssetAmount,
  IPool,
  CompositePool,
  IAssetAmount,
  Amount,
} from "../entities";

import { useField } from "./useField";
import { trimZeros, useBalances } from "./utils";
import { format } from "../utils/format";
export enum SwapState {
  SELECT_TOKENS,
  ZERO_AMOUNTS,
  INSUFFICIENT_FUNDS,
  VALID_INPUT,
  INVALID_AMOUNT,
  INSUFFICIENT_LIQUIDITY,
}

function calculateFormattedPriceImpact(pair: IPool, amount: IAssetAmount) {
  return format(pair.calcPriceImpact(amount), {
    mantissa: 6,
    trimMantissa: true,
  });
}

function calculateFormattedProviderFee(pair: IPool, amount: IAssetAmount) {
  return format(
    pair.calcProviderFee(amount).amount,
    pair.calcProviderFee(amount).asset,
    { mantissa: 5, trimMantissa: true },
  );
}

// TODO: make swap calculator only generate Fractions/Amounts that get stringified in the view
export function useSwapCalculator(input: {
  fromAmount: Ref<string>;
  fromSymbol: Ref<string | null>;
  toAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  balances: Ref<IAssetAmount[]>;
  selectedField: Ref<"from" | "to" | null>;
  slippage: Ref<string>;
  poolFinder: (a: Asset | string, b: Asset | string) => Ref<IPool> | null;
}) {
  // extracting selectedField so we can use it without tracking its change
  let selectedField: "from" | "to" | null = null;
  effect(() => (selectedField = input.selectedField.value));

  // We use a composite pool pair to work out rates
  const pool = computed<IPool | null>(() => {
    if (!input.fromSymbol.value || !input.toSymbol.value) return null;

    if (input.fromSymbol.value === "rowan") {
      return input.poolFinder(input.toSymbol.value, "rowan")?.value ?? null;
    }

    if (input.toSymbol.value === "rowan") {
      return input.poolFinder(input.fromSymbol.value, "rowan")?.value ?? null;
    }

    const fromPair = input.poolFinder(input.fromSymbol.value, "rowan");
    const toPair = input.poolFinder(input.toSymbol.value, "rowan");

    if (!fromPair || !toPair) return null;

    return CompositePool(fromPair.value, toPair.value);
  });

  // Get the balance of the from the users account
  const balance = computed(() => {
    const balanceMap = useBalances(input.balances);
    return input.fromSymbol.value
      ? balanceMap.value.get(input.fromSymbol.value) ?? null
      : null;
  });

  // Get field amounts as domain objects
  const fromField = useField(input.fromAmount, input.fromSymbol);
  const toField = useField(input.toAmount, input.toSymbol);

  // Create a price message eg. 10.123 ATK per BTK
  const priceMessage = computed(() => {
    if (
      !fromField.fieldAmount.value ||
      fromField.fieldAmount.value.equalTo("0") ||
      !pool.value
    ) {
      return "";
    }

    const amount = fromField.fieldAmount.value;

    const pair = pool.value;

    const swapResult = pair.calcSwapResult(amount);

    return `${format(swapResult.divide(amount), {
      mantissa: 6,
      float: true,
    })} ${swapResult.label} per ${amount.label}`;
  });

  // Selected field changes when the user changes the field selection
  // If the selected field is the "tokenA" field and something changes we change the "tokenB" input value
  // If the selected field is the "tokenB" field and something changes we change the "tokenA" input value

  // Changing the "from" field recalculates the "to" amount
  const swapResult = ref<IAssetAmount | null>(null);
  effect(() => {
    if (
      pool.value &&
      fromField.asset.value &&
      fromField.fieldAmount.value &&
      pool.value.contains(fromField.asset.value) &&
      selectedField === "from"
    ) {
      swapResult.value = pool.value.calcSwapResult(fromField.fieldAmount.value);

      const toAmountValue = format(
        swapResult.value.amount,
        swapResult.value.asset,
        {
          mantissa: 10,
          trimMantissa: true,
        },
      );

      input.toAmount.value = toAmountValue;
    }
  });

  // Changing the "to" field recalculates the "from" amount
  const reverseSwapResult = ref<IAssetAmount | null>(null);
  effect(() => {
    if (
      pool.value &&
      toField.asset.value &&
      toField.fieldAmount.value &&
      pool.value.contains(toField.asset.value) &&
      selectedField === "to"
    ) {
      reverseSwapResult.value = pool.value.calcReverseSwapResult(
        toField.fieldAmount.value,
      );

      // Internally trigger calulations based off swapResult as this is how we
      // work out priceImpact, providerFee, minimumReceived

      swapResult.value = pool.value.calcSwapResult(
        reverseSwapResult.value as IAssetAmount,
      );

      input.fromAmount.value = trimZeros(
        format(reverseSwapResult.value.amount, reverseSwapResult.value.asset, {
          mantissa: 8,
        }),
      );
    }
  });

  // Format input amount on blur
  effect(() => {
    if (input.selectedField.value === null && input.toAmount.value) {
      input.toAmount.value = trimZeros(input.toAmount.value);
    }
  });

  // Format input amount on blur
  effect(() => {
    if (input.selectedField.value === null && input.fromAmount.value) {
      input.fromAmount.value = trimZeros(input.fromAmount.value);
    }
  });

  // Cache pool contains asset for reuse as is a little
  const poolContainsFromAsset = computed(() => {
    if (!fromField.asset.value || !pool.value) return false;
    return pool.value.contains(fromField.asset.value);
  });

  const priceImpact = computed(() => {
    if (
      !pool.value ||
      !fromField.asset.value ||
      !fromField.fieldAmount.value ||
      !poolContainsFromAsset.value
    )
      return null;

    return calculateFormattedPriceImpact(
      pool.value as IPool,
      fromField.fieldAmount.value,
    );
  });

  const providerFee = computed(() => {
    if (
      !pool.value ||
      !fromField.asset.value ||
      !fromField.fieldAmount.value ||
      !poolContainsFromAsset.value
    )
      return null;

    return calculateFormattedProviderFee(
      pool.value as IPool,
      fromField.fieldAmount.value,
    );
  });

  // minimumReceived
  const minimumReceived = computed(() => {
    if (!input.slippage.value || !toField.asset.value || !swapResult.value)
      return null;

    const slippage = Amount(input.slippage.value);

    const minAmount = Amount("1")
      .subtract(slippage.divide(Amount("100")))
      .multiply(swapResult.value);

    return AssetAmount(toField.asset.value, minAmount);
  });

  // Derive state
  const state = computed(() => {
    // SwapState.INSUFFICIENT_LIQUIDITY is probably better here
    if (!pool.value) return SwapState.SELECT_TOKENS;
    const fromTokenLiquidity = (pool.value as IPool).amounts.find(
      (amount) => amount.asset.symbol === fromField.asset.value?.symbol,
    );
    const toTokenLiquidity = (pool.value as IPool).amounts.find(
      (amount) => amount.asset.symbol === toField.asset.value?.symbol,
    );

    if (
      !fromTokenLiquidity ||
      !toTokenLiquidity ||
      !fromField.fieldAmount.value ||
      !toField.fieldAmount.value ||
      (fromField.fieldAmount.value?.equalTo("0") &&
        toField.fieldAmount.value?.equalTo("0"))
    ) {
      return SwapState.ZERO_AMOUNTS;
    }

    if (
      toField.fieldAmount.value.greaterThan("0") &&
      fromField.fieldAmount.value.equalTo("0")
    ) {
      return SwapState.INVALID_AMOUNT;
    }

    if (!balance.value?.greaterThanOrEqual(fromField.fieldAmount.value || "0"))
      return SwapState.INSUFFICIENT_FUNDS;

    if (
      fromTokenLiquidity.lessThan(fromField.fieldAmount.value) ||
      toTokenLiquidity.lessThan(toField.fieldAmount.value)
    ) {
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
