import { TransactionStatus } from "ui-core";
import { effect, Ref, ref, ComputedRef } from "@vue/reactivity";
import { IAssetAmount, IAmount, Amount } from "ui-core";

export function getMaxAmount(
  symbol: Ref<string | null>,
  accountBalance: IAssetAmount,
): IAmount {
  if (!symbol) {
    return Amount("0");
  }
  if (symbol.value !== "rowan") {
    return accountBalance;
  } else {
    if (accountBalance.greaterThan(Amount("0.5"))) {
      const fee = 5 * Math.pow(10, accountBalance.decimals - 1); // 0.5 ROWAN
      return accountBalance.subtract(Amount(fee.toString()));
    } else {
      return Amount("0");
    }
  }
}
