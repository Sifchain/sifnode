import { Ref, computed, effect } from "@vue/reactivity";
import {
  Asset,
  AssetAmount,
  IPool,
  CompositePool,
  IAssetAmount,
  Pool,
} from "../entities";
import { useField } from "./useField";
import { assetPriceMessage, trimZeros, useBalances } from "./utils";

export enum SwapState {
  SELECT_TOKENS,
  ZERO_AMOUNTS,
  INSUFFICIENT_FUNDS,
  VALID_INPUT,
}

function calculateFormattedSwapResult(pair: IPool, amount: AssetAmount) {
  return trimZeros(pair.calcSwapResult(amount).toFixed());
}

function calculateFormattedReverseSwapResult(pair: IPool, amount: AssetAmount) {
  return trimZeros(pair.calcReverseSwapResult(amount).toFixed());
}

export function useSwapCalculator(input: {
  fromAmount: Ref<string>;
  fromSymbol: Ref<string | null>;
  toAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  balances: Ref<IAssetAmount[]>;
  selectedField: Ref<"from" | "to" | null>;
  poolFinder: (a: Asset | string, b: Asset | string) => Ref<Pool> | null;
}) {
  // extracting selectedField so we can use it without tracking its change
  let selectedField: "from" | "to" | null = null;
  effect(() => (selectedField = input.selectedField.value));

  // We use a composite pool pair to work out rates
  const pool = computed(() => {
    if (!input.fromSymbol.value || !input.toSymbol.value) return null;

    const fromPair = input.poolFinder(input.fromSymbol.value, "rwn");
    const toPair = input.poolFinder(input.toSymbol.value, "rwn");

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
    const amount = fromField.fieldAmount.value;
    const pair = pool.value;

    return assetPriceMessage(amount, pair, 6);
  });

  // Changing the "from" field recalculates the "to" amount
  effect(() => {
    if (
      pool.value &&
      fromField.asset.value &&
      fromField.fieldAmount.value &&
      selectedField === "from"
    ) {
      input.toAmount.value = calculateFormattedSwapResult(
        pool.value,
        fromField.fieldAmount.value
      );
    }
  });

  // Changing the "to" field recalculates the "from" amount
  effect(() => {
    if (
      pool.value &&
      toField.asset.value &&
      toField.fieldAmount.value &&
      selectedField === "to"
    ) {
      input.fromAmount.value = calculateFormattedReverseSwapResult(
        pool.value,
        toField.fieldAmount.value
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

  // Derive state
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
