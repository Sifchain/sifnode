<template>
  <div class="list">
    <div v-if="walletConnected">
      <table>
        <tr v-for="assetAmount in balances" :key="assetAmount.asset.symbol">
          <td align="left">{{ assetAmount.asset.symbol }}</td>
          <td align="right">{{ assetAmount.toFixed(2) }}</td>
        </tr>
      </table>
    </div>
    <div v-else>No wallet connected</div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { computed } from "@vue/reactivity";

import { useCore } from "../hooks/useCore";
import { Balance } from "../../../core/src";
export default defineComponent({
  name: "ListPage",
  setup() {
    const { store } = useCore();

    const balances = computed(() => {
      // This is trying to simulate mixing real wallet accounts with
      // potential destination swap accounts. It looks the same as the wallet page
      // But eventually will have extra accounts that the wallet doesn't have depending
      // on What is in the top20 tokens
      const balanceSymbols = store.wallet.eth.balances.map(
        (t) => t.asset.symbol
      );
      const tokensNotInBalances = store.asset.top20Tokens.filter((token) => {
        return !balanceSymbols.includes(token.symbol);
      });

      const allBalances = [
        ...store.wallet.eth.balances,
        ...tokensNotInBalances.map((token) => Balance.create(token, "0")),
      ];

      return allBalances;
    });

    const walletConnected = computed(() => store.wallet.eth.isConnected);

    return {
      balances,
      walletConnected,
    };
  },
  components: {},
});
</script>
<style scoped>
table {
  margin: 0 auto;
}
</style>
