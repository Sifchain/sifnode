import { TransactionStatus } from "ui-core";
import { effect, Ref, ref, ComputedRef } from "@vue/reactivity";
import { IAssetAmount, Store } from "ui-core";
import { Fraction } from "../../../../core/src/entities";

export function getMaxAmount(
  symbol: Ref<string | null>,
  accountBalance: IAssetAmount,
): string {
  if (!symbol) {
    return "0";
  }
  if (symbol.value !== "rowan") {
    return accountBalance.toFixed(18);
  } else {
    if (accountBalance.greaterThan(new Fraction("1", "2"))) {
      return accountBalance.subtract(new Fraction("1", "2")).toFixed(18);
    } else {
      return "0";
    }
  }
}
