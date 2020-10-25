import { computed, effect } from "@vue/reactivity";
import { Ref } from "vue";
import { Asset, AssetAmount, IAssetAmount, Pair } from "../../../../core";

function buildAsset(val: string | null) {
  return val === null ? val : Asset.get(val);
}

function buildAssetAmount(asset: Asset | null, amount: string) {
  return asset ? AssetAmount(asset, amount) : asset;
}

export function useField(amount: Ref<string>, symbol: Ref<string | null>) {
  const asset = computed(() => {
    return buildAsset(symbol.value);
  });

  const fieldAmount = computed(() => {
    return buildAssetAmount(asset.value, amount.value);
  });

  return {
    fieldAmount,
    asset,
  };
}

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
  const marketPair = computed(() => {
    if (!input.fromSymbol.value || !input.toSymbol.value) return null;
    return (
      input.marketPairFinder(input.fromSymbol.value, input.toSymbol.value) ??
      null
    );
  });
  const balanceMap = useBalances(input.balances);
  const balance = computed(() => {
    return input.fromSymbol.value
      ? balanceMap.value.get(input.fromSymbol.value) ?? null
      : null;
  });

  const fromField = useField(input.fromAmount, input.fromSymbol);
  const toField = useField(input.toAmount, input.toSymbol);

  const priceAmount = computed(() => {
    if (!fromField.asset.value || !marketPair.value) return null;
    return marketPair.value.priceAsset(fromField.asset.value);
  });

  const nextStepMessage = computed(() => {
    if (!marketPair.value) return "Select tokens";
    if (!balance.value?.greaterThan(fromField.fieldAmount.value || "0"))
      return "Insufficient funds";
    return "Swap";
  });

  effect(() => {
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

    if (
      input.selectedField.value === "from" &&
      marketPair.value &&
      fromField.asset.value &&
      fromField.fieldAmount.value
    ) {
      const asset = fromField.asset.value;
      input.toAmount.value = marketPair.value
        .priceAsset(asset)
        .multiply(fromField.fieldAmount.value)
        .toFixed(asset.decimals);
    }

    if (
      input.selectedField.value === "to" &&
      marketPair.value &&
      toField.asset.value &&
      toField.fieldAmount.value
    ) {
      const asset = toField.asset.value;
      input.fromAmount.value = marketPair.value
        .priceAsset(asset)
        .multiply(toField.fieldAmount.value)
        .toFixed(asset.decimals);
    }
  });

  return {
    priceAmount,
    nextStepMessage,
    toAmount: input.toAmount,
    fromAmount: input.fromAmount,
  };
}
