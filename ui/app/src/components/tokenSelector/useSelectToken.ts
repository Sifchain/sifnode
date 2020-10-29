import { computed } from "@vue/reactivity";
import { Ref } from "vue";
import { Store } from "../../../../core/src";

import { useWallet } from "../../hooks/useWallet";

export function useTokenListing({
  searchText,
  store,
  walletLimit,
  tokenLimit,
}: {
  searchText: Ref<string>;
  store: Store;
  walletLimit: number;
  tokenLimit: number;
}) {
  const { balances } = useWallet(store);

  const walletTokens = computed(() => balances.value.map((tok) => tok.asset));
  const topTokens = computed(() => store.asset.topTokens);
  const fullTokenList = computed(() => {
    return Array.from(new Set([...walletTokens.value, ...topTokens.value]));
  });

  const limitedTokenList = computed(() => {
    return Array.from(
      new Set([
        ...walletTokens.value.slice(0, walletLimit),
        ...topTokens.value.slice(0, tokenLimit),
      ])
    );
  });

  const filteredTokens = computed(() => {
    if (searchText.value) {
      return fullTokenList.value.filter(
        ({ symbol }) =>
          symbol.toLowerCase().indexOf(searchText.value.toLowerCase().trim()) >
          -1
      );
    }
    return limitedTokenList.value;
  });

  return { filteredTokens };
}
