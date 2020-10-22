import { reactive } from "@vue/reactivity";

export function useSwap() {
  const swapState = reactive({
    from: { amount: "0", symbol: null, available: null },
    to: { amount: "0", symbol: null, available: null },
  });
  return { swapState };
}
