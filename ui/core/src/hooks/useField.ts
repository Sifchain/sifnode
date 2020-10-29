import { computed, Ref } from "@vue/reactivity";
import { Asset, AssetAmount, Network, Token } from "../entities";

const TOKENS = {
  atk: Token({
    decimals: 6,
    symbol: "atk",
    name: "AppleToken",
    address: "123",
    network: Network.ETHEREUM,
  }),
  btk: Token({
    decimals: 18,
    symbol: "btk",
    name: "BananaToken",
    address: "1234",
    network: Network.ETHEREUM,
  }),
  eth: Token({
    decimals: 18,
    symbol: "eth",
    name: "Ethereum",
    address: "1234",
    network: Network.ETHEREUM,
  }),
};
function buildAsset(val: string | null) {
  return val === null ? val : Asset.get(val);
}

function buildAssetAmount(asset: Asset | null, amount: string) {
  return asset ? AssetAmount(asset, amount) : asset;
}

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
