import { computed, Ref } from "@vue/reactivity";
import { Amount, Asset, AssetAmount, IAsset } from "../entities";
import { decimalShift } from "../utils/decimalShift";

export function useField(amount: Ref<string>, symbol: Ref<string | null>) {
  const asset = computed(() => {
    if (!symbol.value) return null;
    return Asset(symbol.value);
  });

  const fieldAmount = computed(() => {
    if (!asset.value || !amount.value) return null;
    const shiftedAmountValue = decimalShift(amount.value, asset.value.decimals);
    return AssetAmount(asset.value, shiftedAmountValue);
  });

  return {
    fieldAmount,
    asset,
  };
}
