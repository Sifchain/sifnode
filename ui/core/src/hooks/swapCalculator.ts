import { Ref, computed, effect } from "@vue/reactivity";
import { Asset, AssetAmount, IAssetAmount, Pair } from "../entities";
import { useField } from "./useField";
import { assetPriceMessage, useBalances } from "./utils";

export enum SwapState {
  SELECT_TOKENS,
  ZERO_AMOUNTS,
  INSUFFICIENT_FUNDS,
  VALID_INPUT,
}
function trimZeros(amount: string) {
  return amount.replace(/0+$/, "").replace(/\.$/, ".0");
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

export function useSwapCalculator(input: {
  fromAmount: Ref<string>;
  fromSymbol: Ref<string | null>;
  toAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  balances: Ref<IAssetAmount[]>;
  selectedField: Ref<"from" | "to" | null>;
  marketPairFinder: (a: Asset | string, b: Asset | string) => Pair | null;
}) {
  // We use a market pair to work out the rate
  const marketPair = computed(() => {
    if (!input.fromSymbol.value || !input.toSymbol.value) return null;
    return (
      input.marketPairFinder(input.fromSymbol.value, input.toSymbol.value) ??
      null
    );
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
    const asset = fromField.asset.value;
    const pair = marketPair.value;
    return assetPriceMessage(asset, pair);
  });

  effect(() => {
    input.fromAmount.value = formatValue(
      input.selectedField.value,
      fromField.asset.value,
      input.fromAmount.value
    );

    input.toAmount.value = formatValue(
      input.selectedField.value,
      toField.asset.value,
      input.toAmount.value
    );

    // Changing the "from" field recalculates the "to" amount
    if (
      input.selectedField.value === "from" &&
      marketPair.value &&
      fromField.asset.value &&
      fromField.fieldAmount.value
    ) {
      const asset = fromField.asset.value;
      const assetPrice = marketPair.value.priceAsset(asset);

      input.toAmount.value = assetPrice
        ? trimZeros(
            assetPrice
              .multiply(fromField.fieldAmount.value)
              .toFixed(asset.decimals)
          )
        : "0";
    }

    // Changing the "to" field recalculates the "to" amount
    if (
      input.selectedField.value === "to" &&
      marketPair.value &&
      toField.asset.value &&
      toField.fieldAmount.value
    ) {
      const asset = toField.asset.value;
      const assetPrice = marketPair.value.priceAsset(asset);
      input.fromAmount.value = assetPrice
        ? trimZeros(
            assetPrice
              .multiply(toField.fieldAmount.value)
              .toFixed(asset.decimals)
          )
        : "0";
    }
  });

  const state = computed(() => {
    if (!marketPair.value) return SwapState.SELECT_TOKENS;
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
