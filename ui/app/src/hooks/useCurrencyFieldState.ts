import { effect, Ref, ref } from "@vue/reactivity";

// Store global state between pages
const globalState = {
  fromSymbol: ref<string | null>(null),
  fromAmount: ref<string>("0"),
  toSymbol: ref<string | null>(null),
  toAmount: ref<string>("0"),
  priceImpact: ref<string | null>(null),
  providerFee: ref<string | null>(null),
};

type CurrencyFieldState = {
  fromSymbol: Ref<string | null>;
  fromAmount: Ref<string>;
  toSymbol: Ref<string | null>;
  toAmount: Ref<string>;
  priceImpact: Ref<string | null>;
  providerFee: Ref<string | null>;
};

export function useCurrencyFieldState(): CurrencyFieldState {
  // Copy global state when creating page state
  const fromSymbol = ref<string | null>(globalState.fromSymbol.value);
  const fromAmount = ref<string>(globalState.fromAmount.value);
  const toSymbol = ref<string | null>(globalState.toSymbol.value);
  const toAmount = ref<string>(globalState.toAmount.value);
  const priceImpact = ref<string | null>(globalState.priceImpact.value);
  const providerFee = ref<string | null>(globalState.providerFee.value);

  // Update global state whenchanges occur as sideeffects
  effect(() => (globalState.fromSymbol.value = fromSymbol.value));
  effect(() => (globalState.fromAmount.value = fromAmount.value));
  effect(() => (globalState.toSymbol.value = toSymbol.value));
  effect(() => (globalState.toAmount.value = toAmount.value));
  effect(() => (globalState.priceImpact.value = priceImpact.value));
  effect(() => (globalState.providerFee.value = providerFee.value));

  return {
    fromSymbol,
    fromAmount,
    toSymbol,
    toAmount,
    providerFee,
    priceImpact,
  };
}
