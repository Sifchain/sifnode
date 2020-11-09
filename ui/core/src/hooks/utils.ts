import { computed, Ref } from "@vue/reactivity";
import { AssetAmount, Pool } from "../entities";

export function assetPriceMessage(
  amount: AssetAmount | null,
  pair: Pool | null,
  decimals: number = -1
) {
  if (!pair || !amount || amount.equalTo("0")) return "";
  const swapResult = pair.calcSwapResult(amount);

  const assetPriceStr = [
    swapResult
      .divide(amount)
      .toFixed(decimals > -1 ? decimals : amount.asset.decimals),
    swapResult.asset.symbol.toUpperCase(),
  ].join(" ");

  const formattedPerSymbol = amount.asset.symbol.toUpperCase();

  return `${assetPriceStr} per ${formattedPerSymbol}`;
}

export function trimZeros(amount: string) {
  if (amount.indexOf(".") === -1) return `${amount}.0`;
  return amount.replace(/0+$/, "").replace(/\.$/, ".0");
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
