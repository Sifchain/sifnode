import { Ref, watch } from "@vue/runtime-core";

/**
 * Helper to run effects in the form of a subscription
 * ```ts
 * type Subscription = () => UnsubscribeFn
 * ```
 * @param trigger Trigger the subscriber
 * @param subscriber Synchronous subscription function which returns an unsubscribe function
 *
 */
export function useSubscription(
  sources: Ref<any> | Ref<any>[],
  subscriber: () => () => void,
) {
  watch(sources, (_value, _oldValue, onInvalidateEffect) => {
    const unsubscribe = subscriber();
    onInvalidateEffect(unsubscribe);
  });
}
