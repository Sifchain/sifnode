import { computed, Ref } from "@vue/reactivity";
import { Asset, AssetAmount } from "../entities";
import { toBaseUnits } from "../utils";

export function useField(amount: Ref<string>, symbol: Ref<string | null>) {
  const asset = computed(() => {
    if (!symbol.value) return null;
    return Asset(symbol.value);
  });

  const fieldAmount = computed(() => {
    if (!asset.value || !amount.value) return null;
    return AssetAmount(asset.value, toBaseUnits(amount.value, asset.value));
  });

  return {
    fieldAmount,
    asset,
  };
}
