import { Ref, computed, effect } from "@vue/reactivity";
import { Asset, AssetAmount, IAssetAmount, Pair } from "../entities";
import { useField } from "./useField";
import { assetPriceMessage } from "./utils";

function useBalances(balances: Ref<AssetAmount[]>) {
  return computed(() => {
    const map = new Map<string, AssetAmount>();

    for (const item of balances.value) {
      map.set(item.asset.symbol, item);
    }
    return map;
  });
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

  const nextStepMessage = computed(() => {
    if (!marketPair.value) return "Select tokens";
    if (
      fromField.fieldAmount.value?.equalTo("0") &&
      toField.fieldAmount.value?.equalTo("0")
    )
      return "You cannot swap 0";
    if (!balance.value?.greaterThan(fromField.fieldAmount.value || "0"))
      return "Insufficient funds";
    return "Swap";
  });

  const canSwap = computed(() => {
    return nextStepMessage.value === "Swap";
  });

  effect(() => {
    // Deselect a field formats all values
    if (input.selectedField.value === null) {
      const fromAsset = fromField.asset.value;
      if (fromAsset) {
        input.fromAmount.value = AssetAmount(
          fromAsset,
          input.fromAmount.value
        ).toFixed();
      }

      const toAsset = fromField.asset.value;
      if (toAsset) {
        input.toAmount.value = AssetAmount(
          toAsset,
          input.toAmount.value
        ).toFixed();
      }
    }

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
        ? assetPrice
            .multiply(fromField.fieldAmount.value)
            .toFixed(asset.decimals)
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
        ? assetPrice.multiply(toField.fieldAmount.value).toFixed(asset.decimals)
        : "0";
    }
  });

  return {
    canSwap,
    priceMessage,
    nextStepMessage,
    fromFieldAmount: fromField.fieldAmount,
    toFieldAmount: toField.fieldAmount,
    toAmount: input.toAmount,
    fromAmount: input.fromAmount,
  };
}
