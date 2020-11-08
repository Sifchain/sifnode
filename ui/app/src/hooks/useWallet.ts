import { computed, ComputedRef } from "@vue/reactivity";
import { IAssetAmount, Store } from "../../../core/src";

export function useWallet(
  store: Store
): { balances: ComputedRef<IAssetAmount[]> } {
  const balances = computed(() => store.wallet.sif.balances);

  return { balances };
}
