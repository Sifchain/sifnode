import { computed, Ref } from "@vue/reactivity";
import { Asset, AssetAmount } from "../entities";

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
