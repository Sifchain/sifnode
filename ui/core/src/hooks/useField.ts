import { computed, Ref } from "@vue/reactivity";
import { Amount, Asset, AssetAmount, IAsset } from "../entities";

export function useField(amount: Ref<string>, symbol: Ref<string | null>) {
  const asset = computed(() => {
    if (!symbol.value) return null;
    return Asset(symbol.value);
  });

  const fieldAmount = computed(() => {
    if (!asset.value) return null;
    return AssetAmount(asset.value, Amount(amount.value));
  });

  return {
    fieldAmount,
    asset,
  };
}
