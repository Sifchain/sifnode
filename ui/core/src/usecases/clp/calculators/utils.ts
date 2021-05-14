import { computed, Ref } from "@vue/reactivity";
import { IAssetAmount, IPool } from "../../../entities";
import { format } from "../../../utils/format";

export function assetPriceMessage(
  amount: IAssetAmount | null,
  pair: IPool | null,
  decimals: number,
) {
  if (!pair || !amount || amount.equalTo("0")) return "";
  const swapResult = pair.calcSwapResult(amount);

  return `${format(swapResult.divide(amount), {
    mantissa: decimals,
  })} ${swapResult.label} per ${amount.label}`;
}

export function trimZeros(amount: string) {
  if (amount.indexOf(".") === -1) return `${amount}.0`;
  const tenDecimalsMax = parseFloat(amount).toFixed(10);
  return tenDecimalsMax.replace(/0+$/, "").replace(/\.$/, ".0");
}

export function useBalances(balances: Ref<IAssetAmount[]>) {
  return computed(() => {
    const map = new Map<string, IAssetAmount>();

    for (const item of balances.value) {
      map.set(item.asset.symbol, item);
    }
    return map;
  });
}
