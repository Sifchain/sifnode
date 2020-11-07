import {
  Ref,
  unref,
  computed,
  effect,
  // pauseTracking,
  // resetTracking,
  enableTracking,
} from "@vue/reactivity";
import { watch } from "vue";
import {
  Asset,
  AssetAmount,
  CompositePair,
  IAssetAmount,
  Pair,
} from "../entities";
import { useField } from "./useField";
import { assetPriceMessage, trimZeros, useBalances } from "./utils";

export enum SwapState {
  SELECT_TOKENS,
  ZERO_AMOUNTS,
  INSUFFICIENT_FUNDS,
  VALID_INPUT,
}

function formatValue(
  selectedField: string | null,
  asset: Asset | null,
  amount: string
) {
  if (selectedField === null) {
    if (asset) {
      return trimZeros(AssetAmount(asset, amount).toFixed());
    }
  }
  return amount;
}

function calculateSwapResult(pair: Pair, amount: AssetAmount) {
  return trimZeros(pair.calcSwapResult(amount).toFixed());
}

function calculateReverseSwapResult(pair: Pair, amount: AssetAmount) {
  return trimZeros(pair.calcReverseSwapResult(amount).toFixed());
}

export function useSwapCalculator(input: {
  fromAmount: Ref<string>;
  fromSymbol: Ref<string | null>;
  toAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  balances: Ref<IAssetAmount[]>;
  selectedField: Ref<"from" | "to" | null>;
  marketPairFinder: (a: Asset | string, b: Asset | string) => Pair | null;
}) {
  // extracting selectedField so we can use it without tracking its change
  let selectedField: "from" | "to" | null = null;
  effect(() => (selectedField = input.selectedField.value));

  // We use a composite market pair to work out rates
  const pool = computed(() => {
    if (!input.fromSymbol.value || !input.toSymbol.value) return null;

    const fromPair = input.marketPairFinder(input.fromSymbol.value, "rwn");
    const toPair = input.marketPairFinder(input.toSymbol.value, "rwn");

    if (!fromPair || !toPair) return null;

    return CompositePair(fromPair, toPair);
  });

  // get the balance of the from account
  const balance = computed(() => {
    const balanceMap = useBalances(input.balances);
    return input.fromSymbol.value
      ? balanceMap.value.get(input.fromSymbol.value) ?? null
      : null;
  });

  // Get field amounts as domain objects
  const fromField = useField(input.fromAmount, input.fromSymbol);
  const toField = useField(input.toAmount, input.toSymbol);

  // Create a price message
  const priceMessage = computed(() => {
    const amount = fromField.fieldAmount.value;
    const pair = pool.value;

    return assetPriceMessage(amount, pair, 6);
  });

  effect(() => {
    // Changing the "from" field recalculates the "to" amount
    if (
      pool.value &&
      fromField.asset.value &&
      fromField.fieldAmount.value &&
      selectedField === "from"
    ) {
      input.toAmount.value = calculateSwapResult(
        pool.value,
        fromField.fieldAmount.value
      );
    }
  });

  effect(() => {
    // Changing the "to" field recalculates the "from" amount
    if (
      pool.value &&
      toField.asset.value &&
      toField.fieldAmount.value &&
      selectedField === "to"
    ) {
      input.fromAmount.value = calculateReverseSwapResult(
        pool.value,
        toField.fieldAmount.value
      );
    }
  });
  effect(() => {
    if (input.selectedField.value === null && input.toAmount.value) {
      input.toAmount.value = trimZeros(input.toAmount.value);
    }
  });
  effect(() => {
    if (input.selectedField.value === null && input.fromAmount.value) {
      input.fromAmount.value = trimZeros(input.fromAmount.value);
    }
  });
  const state = computed(() => {
    if (!pool.value) return SwapState.SELECT_TOKENS;
    if (
      fromField.fieldAmount.value?.equalTo("0") &&
      toField.fieldAmount.value?.equalTo("0")
    )
      return SwapState.ZERO_AMOUNTS;
    if (!balance.value?.greaterThan(fromField.fieldAmount.value || "0"))
      return SwapState.INSUFFICIENT_FUNDS;

    return SwapState.VALID_INPUT;
  });

  return {
    priceMessage,
    state,
    fromFieldAmount: fromField.fieldAmount,
    toFieldAmount: toField.fieldAmount,
    toAmount: input.toAmount,
    fromAmount: input.fromAmount,
  };
}
