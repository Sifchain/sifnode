import { Ref, ref } from "@vue/reactivity";

const fromSymbol = ref<string | null>(null);
const fromAmount = ref<string>("0");
const toSymbol = ref<string | null>(null);
const toAmount = ref<string>("0");

type CurrencyFieldState = {
  fromSymbol: Ref<string | null>;
  fromAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  toAmount: Ref<string>;
};

export function useCurrencyFieldState(): CurrencyFieldState {
  return {
    fromSymbol,
    fromAmount,
    toSymbol,
    toAmount,
  };
}
