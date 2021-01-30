import { computed, Ref } from "@vue/reactivity";
import ColorHash from "color-hash";
import { Asset, Network } from "ui-core";

export function formatSymbol(symbol: string) {
  if (symbol.indexOf("c") === 0) {
    return ["c", symbol.slice(1).toUpperCase()].join("");
  }
  return symbol.toUpperCase();
}

// TODO: These could be replaced with a look up table
export function getPeggedSymbol(symbol: string) {
  if (symbol === "erowan") return "rowan";
  return "c" + symbol;
}
export function getUnpeggedSymbol(symbol: string) {
  if (symbol === "rowan") return "erowan";
  return symbol.indexOf("c") === 0 ? symbol.slice(1) : symbol;
}

export function getAssetLabel(t: Asset) {
  if (t.network === Network.SIFCHAIN) {
    return formatSymbol(t.symbol);
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
