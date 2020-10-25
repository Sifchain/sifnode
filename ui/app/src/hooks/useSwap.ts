import { ref } from "@vue/reactivity";

const fromAmount = ref("0");
const fromSymbol = ref<string | null>(null);
const toAmount = ref("0");
const toSymbol = ref<string | null>(null);

export function useSwap() {
  return {
    from: { amount: fromAmount, symbol: fromSymbol },
    to: { amount: toAmount, symbol: toSymbol },
  };
}
