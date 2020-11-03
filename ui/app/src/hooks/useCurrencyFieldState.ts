import { effect, Ref, ref } from "@vue/reactivity";

// Store global state between pages
const globalState = {
  fromSymbol: ref<string | null>(null),
  fromAmount: ref<string>("0"),
  toSymbol: ref<string | null>(null),
  toAmount: ref<string>("0"),
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
  const fromAmount = ref<string>(globalState.fromAmount.value);
  const toSymbol = ref<string | null>(globalState.toSymbol.value);
  const toAmount = ref<string>(globalState.toAmount.value);

  // Update global state whenchanges occur as sideeffects
  effect(() => (globalState.fromSymbol.value = fromSymbol.value));
  effect(() => (globalState.fromAmount.value = fromAmount.value));
  effect(() => (globalState.toSymbol.value = toSymbol.value));
  effect(() => (globalState.toAmount.value = toAmount.value));

  return {
    fromSymbol,
    fromAmount,
    toSymbol,
    toAmount,
  };
}
