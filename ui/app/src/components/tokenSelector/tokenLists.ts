import { Ref, computed, ComputedRef, ref, effect } from "@vue/reactivity";
import { Asset, Store } from "ui-core";

export function generateTokenSearchLists({
  walletLimit = 10,
  tokenLimit = 20,
  walletTokens = [],
  topTokens = [],
}: {
  walletLimit?: number;
  tokenLimit?: number;
  walletTokens?: Asset[];
  topTokens?: Asset[];
}) {
  // You can search through a larger list than we display by default
  // as we need to show tokens a user will have in their wallet

  // in order to comply with specification in the PRD we need to
  // generate a combination of lists based on limits provided
  // fullTokenList for searching = top tokens and wallet tokens
  const fullSearchList = computed(() => {
    return Array.from(new Set([...walletTokens, ...topTokens]));
  });

  // List for default display as we have limits on the wallet
  const displayList = computed(() => {
    return Array.from(
      new Set([
        ...walletTokens.slice(0, walletLimit),
        ...topTokens.slice(0, tokenLimit),
      ]),
    );
  });

  return { fullSearchList, displayList };
}

export function filterTokenList({
  searchText,
  tokens,
  displayList = ref([]),
}: {
  searchText: Ref<string>;
  tokens: Ref<Asset[]>;
  displayList?: Ref<Asset[]>;
}) {
  return computed(() => {
    console.log({ tokens: tokens.value });

    const list = searchText.value
      ? tokens.value.filter(
          ({ symbol }) =>
            symbol
              .toLowerCase()
              .indexOf(searchText.value.toLowerCase().trim()) > -1,
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
  return computed(
    () =>
      list.value?.map(item =>
        selectedTokens.includes(item.symbol)
          ? { disabled: true, ...item }
          : item,
      ) ?? [],
  );
}
