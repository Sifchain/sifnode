<script lang="ts">
import { defineComponent, PropType } from "vue"; /* eslint-disable-line */
import { ref } from "@vue/reactivity";
import { useCore } from "../../hooks/useCore";
import AssetItem from "@/components/shared/AssetItem.vue";
import { useTokenListing } from "./useSelectToken";
import SifInput from "@/components/shared/SifInput.vue";

export default defineComponent({
  name: "SelectTokenDialog",
  components: { AssetItem, SifInput },
  emits: ["tokenselected"],
  props: { selectedTokens: Array as PropType<string[]> },
  setup(props, context) {
    const { store, actions } = useCore();

    const searchText = ref("");

    const { filteredTokens } = useTokenListing({
      searchText,
      store,
      tokenLimit: 20,
      walletLimit: 10,
      selectedTokens: props.selectedTokens || [],
      // topTokens: actions.p
    });

    function selectToken(symbol: string) {
      context.emit("tokenselected", symbol);
    }

    return { filteredTokens, searchText, selectToken };
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
    <div class="no-tokens-message" v-if="filteredTokens.length === 0">
      <p>No tokens available.</p>
    </div>
    <button
      class="option"
      v-for="token in filteredTokens"
      :disabled="token.disabled"
      :key="token.symbol"
      @click="selectToken(token.symbol)"
    >
      <AssetItem :symbol="token.symbol" />
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