import { computed } from "@vue/reactivity";
import { useCore } from "./useCore";
import { useSubscription } from "./useSubscrition";

export function useInitialize() {
  const { actions, store } = useCore();
  // initialize subscriptions
  useSubscription(
    computed(() => store.wallet.eth.address), // Needs a ref
    actions.peg.subscribeToUnconfirmedPegTxs
  );
}
