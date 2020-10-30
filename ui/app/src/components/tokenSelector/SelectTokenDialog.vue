<script lang="ts">
import { defineComponent, PropType } from "vue";
import { ref } from "@vue/reactivity";
import { useCore } from "../../hooks/useCore";
import AssetItem from "./AssetItem.vue";
import { useTokenListing } from "./useSelectToken";

export default defineComponent({
  name: "SelectTokenDialog",
  components: { AssetItem },
  emits: ["tokenselected"],
  props: { selectedTokens: Array as PropType<string[]> },
  setup(props, context) {
    const { store } = useCore();

    const searchText = ref("");

    const { filteredTokens } = useTokenListing({
      searchText,
      store,
      tokenLimit: 20,
      walletLimit: 10,
      selectedTokens: props.selectedTokens || [],
    });

    function selectToken(symbol: string) {
      context.emit("tokenselected", symbol);
    }

    return { filteredTokens, searchText, selectToken };
  },
});
</script>

<template>
  <p>Select a token</p>
  <input class="search-input" v-model="searchText" />
  <p>Token Name</p>
  <hr />
  <div class="token-list">
    <button
      class="token-button"
      v-for="token in filteredTokens"
      :disabled="token.disabled"
      :key="token.symbol"
      @click="selectToken(token.symbol)"
    >
      <AssetItem :symbol="token.symbol" />
    </button>
  </div>
</template>

<style scoped>
.token-list {
  display: flex;
  flex-direction: column;
  max-height: 50vh;
  overflow-y: auto;
}
.token-button {
  text-align: left;
  background: transparent;
  border: none;
  margin-bottom: 0.5rem;
}
</style>