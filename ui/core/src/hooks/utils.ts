import { computed, Ref } from "@vue/reactivity";
import { Asset, Token, AssetAmount, IPool, Pool } from "../entities";

export function assetPriceMessage(
  amount: AssetAmount | null,
  pair: IPool | null,
  decimals: number = -1
) {
  if (!pair || !amount || amount.equalTo("0")) return "";
  const swapResult = pair.calcSwapResult(amount);

  const assetPriceStr = [
    swapResult
      .divide(amount)
      .toFixed(decimals > -1 ? decimals : amount.asset.decimals),
    swapResult.asset.symbol.toLowerCase().includes("rowan")
      ? swapResult.asset.symbol.toUpperCase()
      : "c" + swapResult.asset.symbol.slice(1).toUpperCase(),
  ].join(" ");

  const formattedPerSymbol = amount.asset.symbol.toLowerCase().includes("rowan")
    ? amount.asset.symbol.toUpperCase()
    : "c" + amount.asset.symbol.slice(1).toUpperCase();

  return `${assetPriceStr} per ${formattedPerSymbol}`;
}

export function trimZeros(amount: string) {
  if (amount.indexOf(".") === -1) return `${amount}.0`;
  const tenDecimalsMax = parseFloat(amount).toFixed(10);
  return tenDecimalsMax.replace(/0+$/, "").replace(/\.$/, ".0");
}

export function useBalances(balances: Ref<AssetAmount[]>) {
  return computed(() => {
    const map = new Map<string, AssetAmount>();

    for (const item of balances.value) {
      map.set(item.asset.symbol, item);
    }
    return map;
  });
}

export function buildAsset(val: string | null) {
  return val === null ? val : Asset.get(val);
}

export function buildAssetAmount(asset: Asset | null, amount: string) {
  return asset ? AssetAmount(asset, amount) : asset;
}
