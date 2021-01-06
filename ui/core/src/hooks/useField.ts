import { computed, Ref } from "@vue/reactivity";
import { buildAsset, buildAssetAmount } from "./utils";

export function useField(amount: Ref<string>, symbol: Ref<string | null>) {
  const asset = computed(() => {
    if (!symbol.value) return null;
    return buildAsset(symbol.value);
  });

  const fieldAmount = computed(() => {
    if (!asset.value) return null;
    return buildAssetAmount(asset.value, amount.value);
  });

  return {
    fieldAmount,
    asset,
  };
}
