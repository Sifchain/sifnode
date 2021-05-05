import { TransactionStatus } from "ui-core";
import { effect, Ref, ref, ComputedRef } from "@vue/reactivity";
import { IAssetAmount, IAmount, Amount } from "ui-core";

// We set this static fee to minus from some ROWAN transactions such
// that users don't have to manually minus it from KEPLR
const ROWAN_GAS_FEE = Amount("500000000000000000"); // 0.5 ROWAN

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
    if (accountBalance.greaterThan(ROWAN_GAS_FEE)) {
      const fee = 5 * Math.pow(10, accountBalance.decimals - 1);
      return accountBalance.subtract(Amount(fee.toString()));
    } else {
      return Amount("0");
    }
  }
}
