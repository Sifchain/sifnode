<script lang="ts">
import { defineComponent, PropType } from "vue"; /* eslint-disable-line */
import { Ref, ref, toRefs } from "@vue/reactivity";
import { useCore } from "../../hooks/useCore";
import AssetItem from "@/components/shared/AssetItem.vue";
import {
  disableSelected,
  filterTokenList,
  generateTokenSearchLists,
} from "./tokenLists";
import SifInput from "@/components/shared/SifInput.vue";
import { Asset } from "ui-core";

type Balances = {
  age: number;
};

export default defineComponent({
  name: "SelectTokenDialog",
  components: { AssetItem, SifInput },
  emits: ["tokenselected"],
  props: {
    selectedTokens: Array as PropType<string[]>,
    forceShowAllATokens: { type: Boolean, default: false },
    displayList: { type: Object as PropType<Asset[]>, default: ref([]) },
    fullSearchList: {
      type: Object as PropType<Asset[]>,
      default: ref([]),
    },
  },
  setup(props, context) {
    const { store, actions } = useCore();
    const { forceShowAllATokens } = props;
    const searchText = ref("");
    const selectedTokens = props.selectedTokens || [];
    const allSifTokens = ref(actions.peg.getSifTokens());
    const { fullSearchList, displayList } = toRefs(props);

    const list = filterTokenList({
      searchText,
      tokens: true ? allSifTokens : fullSearchList,
      displayList: forceShowAllATokens ? allSifTokens : displayList,
    });
    console.log("allSifTokens", allSifTokens);
    const balances = store.wallet.sif.balances;
    // const balancesObj: Balances = {};
    // const assetBalances = balances.forEach((balance) => {
    //   console.log("hello");
    //   console.log(balance.asset.symbol);
    //   console.log(balance.toFixed());
    //   balancesObj[balance.asset.symbol] = balance.toFixed();
    // });
    // console.log("balances1", balances);
    // console.log("balances", balances.value);
    let tokens = disableSelected({ list, selectedTokens });

    function selectToken(symbol: string) {
      context.emit("tokenselected", symbol);
    }
    console.log("toaaa", tokens);
    // const sifBalances =
    // const tokenBalances = tokens.value.map((token) => {
    // return { hey: 1 };
    // balances.find((balance) => {
    //   return balance.asset.symbol.toLowerCase() === symbol.value.toLowerCase();
    // });
    // });
    return { balances, tokens, searchText, selectToken };
  },
});
</script>

<template>
  <div class="header">
    <h3 class="title">Select a token</h3>
    <SifInput
      gold
      placeholder="Search name or paste address"
      class="sif-input"
      type="text"
      v-model="searchText"
    />
    <h4 class="list-title">Token Name</h4>
  </div>

  <div class="body">
    <div class="no-tokens-message" v-if="tokens.length === 0">
      <p>No tokens available.</p>
    </div>
    <button
      class="option"
      v-for="token in tokens"
      :disabled="token.disabled"
      :key="token.symbol"
      @click="selectToken(token.symbol)"
    >
      <AssetItem :symbol="token.symbol" />
      <div>
        {{
          balances
            .find((balance) => {
              return (
                balance.asset.symbol.toLowerCase() ===
                token.symbol.toLowerCase()
              );
            })
            .toFixed()
        }}
      </div>
    </button>
  </div>
</template>

<style lang="scss" scoped>
.token-list {
  display: flex;
  flex-direction: column;
  max-height: 50vh;
  overflow-y: auto;
}

.title {
  font-size: $fs_lg;
  color: $c_text;
  margin-bottom: 1rem;
  text-align: left;
}
.list-title {
  color: $c_text;
  text-align: left;
  margin-top: 30px;
  margin-bottom: 1rem;
}

.header {
  position: relative;
  padding: 40px 15px 0;
}

.body {
  padding-top: 14px;
  flex-grow: 1;
  max-height: 50vh;
  overflow-y: scroll;
  display: flex;
  flex-direction: column;
  border-top: $divider;
}

.option {
  margin-bottom: 22px;
  font-size: $fs_md;
  font-weight: bold;
  text-align: left;
  color: $c_text;
  padding-left: 15px;
  cursor: pointer;
  text-align: left;
  background: transparent;
  border: none;
  display: flex;
  justify-content: space-between;
  @include listAnimation;

  &[disabled] {
    color: #bbb;
    pointer-events: none;
  }
}
.no-tokens-message {
  padding: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
