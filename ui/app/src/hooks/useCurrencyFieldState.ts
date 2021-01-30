import { effect, Ref, ref } from "@vue/reactivity";

// Store global state between pages
const globalState = {
  fromSymbol: ref<string | null>(null),
  toSymbol: ref<string | null>(null),
};

type CurrencyFieldState = {
  fromSymbol: Ref<string | null>;
  fromAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  toAmount: Ref<string>;
};

export function useCurrencyFieldState(): CurrencyFieldState {
  // Copy global state when creating page state
  const fromSymbol = ref<string | null>(globalState.fromSymbol.value);
  const toSymbol = ref<string | null>(globalState.toSymbol.value);

  // Local page state
  const fromAmount = ref<string>("0");
  const toAmount = ref<string>("0");


  // Update global state whenchanges occur as sideeffects
  effect(() => (globalState.fromSymbol.value = fromSymbol.value));
  effect(() => (globalState.toSymbol.value = toSymbol.value));

  return {
    fromSymbol,
    fromAmount,
    toSymbol,
    toAmount,
  };
}
