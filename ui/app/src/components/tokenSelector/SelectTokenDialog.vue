<script lang="ts">
import { defineComponent } from "vue";
import { ref } from "@vue/reactivity";
import { useCore } from "../../hooks/useCore";
import AssetItem from "./AssetItem.vue";
import { useSwap } from "@/hooks/useSwap";
import { useTokenListing } from "./useSelectToken";

export default defineComponent({
  name: "SelectTokenDialog",
  emits: ["close"],
  components: { AssetItem },
  props: { label: String },
  setup(props, context) {
    const { store } = useCore();
    const swapState = useSwap();

    const searchText = ref("");

    const { filteredTokens } = useTokenListing({
      searchText,
      store,
      tokenLimit: 20,
      walletLimit: 10,
    });

    function selectToken(symbol: string) {
      const label = props.label?.toLowerCase() as "from" | "to";

      swapState[label].symbol.value = symbol;

      context.emit("close");
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
    <div
      class="token-button"
      v-for="token in filteredTokens"
      :key="token.symbol"
      @click="selectToken(token.symbol)"
    >
      <AssetItem :symbol="token.symbol" />
    </div>
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