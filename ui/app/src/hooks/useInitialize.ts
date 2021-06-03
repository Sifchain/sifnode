import { computed } from "@vue/reactivity";
import { useCore } from "./useCore";
import { useSubscription } from "./useSubscrition";

export function useInitialize() {
  const { usecases, store } = useCore();
  // initialize subscriptions
  useSubscription(
    computed(() => store.wallet.eth.address), // Needs a ref
    usecases.peg.subscribeToUnconfirmedPegTxs,
  );
}
