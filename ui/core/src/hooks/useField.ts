import { computed, Ref } from "@vue/reactivity";
import { Asset, AssetAmount, Network, Token } from "../entities";
import { buildAsset, buildAssetAmount } from "./utils";

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
