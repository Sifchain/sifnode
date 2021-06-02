<script lang="ts">
import { defineComponent, PropType } from "vue"; /* eslint-disable-line */
import { Ref, ref, toRefs } from "@vue/reactivity";
import { useCore } from "../../hooks/useCore";
import { sortAssetAmount } from "../../views/utils/sortAssetAmount";
import AssetItem from "@/components/shared/AssetItem.vue";
import { format } from "ui-core/src/utils/format";
import { computed } from "@vue/reactivity";
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
    const { store, usecases } = useCore();
    const { forceShowAllATokens } = props;
    const searchText = ref("");
    const selectedTokens = props.selectedTokens || [];
    const allSifTokens = ref(usecases.peg.getSifTokens());
    const { fullSearchList, displayList } = toRefs(props);

    function selectToken(symbol: string) {
      context.emit("tokenselected", symbol);
    }

    const tokenBalances = computed(() => {
      const list = filterTokenList({
        searchText,
        tokens: true ? allSifTokens : fullSearchList,
        displayList: forceShowAllATokens ? allSifTokens : displayList,
      });

      const balances = store.wallet.sif.balances;

      let tokens = disableSelected({ list, selectedTokens });

      let tokenBalances = tokens.value.map((asset) => {
        let balance = null;
        // If not connected to Keplr, we still want to display the possible assets to trade
        if (balances) {
          balance = balances.find((balance) => {
            return (
              balance.asset.symbol.toLowerCase() === asset.symbol.toLowerCase()
            );
          });
        }
        return { asset: asset, amount: balance };
      });

      tokenBalances = sortAssetAmount(tokenBalances);

      return tokenBalances;
    });
    return { searchText, selectToken, tokenBalances, format };
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
    <div class="no-tokens-message" v-if="tokenBalances.length === 0">
      <p>No tokens available.</p>
    </div>
    <button
      :data-handle="tb.asset.symbol + '-select-button'"
      class="option"
      v-for="tb in tokenBalances"
      :disabled="tb.asset.disabled"
      :key="tb.asset.symbol"
      @click="selectToken(tb.asset.symbol)"
    >
      <AssetItem :symbol="tb.asset.symbol" />
      <div class="balance">
        {{ tb.amount ? format(tb.amount, tb.asset, { mantissa: 4 }) : "0" }}
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
  color: #818181;
  font-family: "PT Serif", serif;
  font-size: 14px;
  font-style: italic;
  padding-right: 14px;
  &[disabled] {
    color: #bbb;
    pointer-events: none;
  }
}
.balance {
  color: #818181;
  font-family: "PT Serif", serif;
  font-size: 14px;
  font-style: italic;
}
.no-tokens-message {
  padding: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
