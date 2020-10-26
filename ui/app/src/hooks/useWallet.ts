import { computed } from "@vue/reactivity";
import { Store } from "../../../core/src";

export function useWallet(store: Store) {
  const balances = computed(() => [
    ...store.wallet.eth.balances,
    ...store.wallet.sif.balances,
  ]);

  return { balances };
}
