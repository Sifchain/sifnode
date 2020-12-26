import { Ref, computed, ComputedRef } from "@vue/reactivity";
import { Asset, Store } from "ui-core";
import { useWallet } from "../../hooks/useWallet";

export function generateTokenSearchLists({
  store,
  walletLimit,
  tokenLimit,
  topTokens = [],
}: {
  store: Store;
  walletLimit: number;
  tokenLimit: number;
  topTokens?: Asset[];
}) {
  const { balances } = useWallet(store);

  const walletTokens = computed(() => balances.value.map((tok) => tok.asset));
  // You can search through a larger list than we display by default
  // as we need to show tokens a user will have in their wallet

  // in order to comply with specification in the PRD we need to
  // generate a combination of lists based on limits provided
  // fullTokenList for searching = top tokens and wallet tokens
  const fullSearchList = computed(() => {
    return Array.from(new Set([...walletTokens.value, ...topTokens]));
  });

  // List for default display as we have limits on the wallet
  const displayList = computed(() => {
    return Array.from(
      new Set([
        ...walletTokens.value.slice(0, walletLimit),
        ...topTokens.slice(0, tokenLimit),
      ])
    );
  });
  return { fullSearchList, displayList };
}

export function filterTokenList({
  searchText,
  tokens,
  displayList,
}: {
  searchText: Ref<string>;
  tokens: Ref<Asset[]>;
  displayList?: Ref<Asset[]>;
}) {
  return computed(() => {
    const list = searchText.value
      ? tokens.value.filter(
          ({ symbol }) =>
            symbol
              .toLowerCase()
              .indexOf(searchText.value.toLowerCase().trim()) > -1
        )
      : (displayList || tokens).value;

    return list;
  });
}

export function disableSelected({
  list,
  selectedTokens = [],
}: {
  list: Ref<Asset[]>;
  selectedTokens: string[];
}) {
  return computed(() =>
    list.value.map((item) =>
      selectedTokens.includes(item.symbol) ? { disabled: true, ...item } : item
    )
  );
}

export function useTokenListing({
  searchText,
  store,
  walletLimit,
  tokenLimit,
  selectedTokens = [],
  topTokens = [],
}: {
  searchText: Ref<string>;
  store: Store;
  walletLimit: number;
  tokenLimit: number;
  selectedTokens: string[];
  topTokens?: Asset[];
}): { filteredTokens: ComputedRef<Asset[]> } {
  const { displayList, fullSearchList } = generateTokenSearchLists({
    store,
    tokenLimit,
    walletLimit,
    topTokens,
  });
  const list = filterTokenList({
    searchText,
    tokens: fullSearchList,
    displayList,
  });
  const filteredTokens = disableSelected({ list, selectedTokens });
  return { filteredTokens };
}
