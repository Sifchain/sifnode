import { computed, Ref } from "@vue/reactivity";
import ColorHash from "color-hash";
import { Asset, Network, TxHash } from "ui-core";

export function formatSymbol(symbol: string) {
  if (symbol.indexOf("c") === 0) {
    return ["c", symbol.slice(1).toUpperCase()].join("");
  }
  return symbol.toUpperCase();
}

export function formatPercentage(amount: string) {
  return parseFloat(amount) < 0.01
    ? "< 0.01%"
    : `${parseFloat(amount).toFixed(2)}%`;
}
// TODO: make this work for AssetAmounts and Fractions / Amounts
export function formatNumber(displayNumber: string) {
  if (!displayNumber) return "0"
  const amount = parseFloat(displayNumber);
  if (amount < 100000) {
    return amount.toFixed(5);
  } else {
    return amount.toFixed(2);
  }
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

  return {
    token: token,
    label: tokenLabel,
    background: backgroundStyle,
  };
}

export function getBlockExplorerUrl(chainId: string, txHash?: TxHash): string {
  if (chainId === "sifchain") {
    if(!txHash) return "https://blockexplorer.sifchain.finance/";
    return `https://blockexplorer.sifchain.finance/transactions/${txHash}`;
  } else {
    if (!txHash) return `https://blockexplorer-${chainId}.sifchain.finance/`;
    return `https://blockexplorer-${chainId}.sifchain.finance/transactions/${txHash}`;
  }
}
