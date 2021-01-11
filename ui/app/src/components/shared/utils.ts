import { computed, Ref } from "@vue/reactivity";
import ColorHash from "color-hash";
import { Asset, Network } from "ui-core";

export function getAssetLabel(t: Asset) {
  if (t.network === Network.SIFCHAIN && t.symbol.indexOf("c") === 0) {
    return ["c", t.symbol.slice(1).toUpperCase()].join("");
  }

  if (t.network === Network.ETHEREUM && t.symbol.toLowerCase() === "erowan") {
    return "eROWAN";
  }
  return t.symbol.toUpperCase();
}

export function useAssetItem(symbol: Ref<string | undefined>) {
  const token = computed(() =>
    symbol.value ? Asset.get(symbol.value) : undefined
  );

  const tokenLabel = computed(() => {
    if (!token.value) return "";
    return getAssetLabel(token.value);
  });

  const backgroundStyle = computed(() => {
    if (!symbol.value) return "";

    const colorHash = new ColorHash();

    const color = symbol ? colorHash.hex(symbol.value) : [];

    return `background: ${color};`;
  });

  const asset = {
    token: token,
    label: tokenLabel,
    background: backgroundStyle,
  };

  return asset;
}
