<template>
  <div class="list">
    <table>
      <tr v-for="assetAmount in balances" :key="assetAmount.asset.symbol">
        <td align="right">{{ assetAmount.toFixed(6) }}</td>
        <td align="left">{{ assetAmount.asset.symbol }}</td>
      </tr>
    </table>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { computed } from "@vue/reactivity";

import { useCore } from "../hooks/useCore";
import { Balance } from "../../../core/src";
export default defineComponent({
  name: "ListPage",
  setup() {
    const { store, actions } = useCore();

    const balances = computed(() => {
      const balanceSymbols = store.wallet.balances.map((t) => t.asset.symbol);
      const tokensNotInBalances = store.asset.top20Tokens.filter((token) => {
        return !balanceSymbols.includes(token.symbol);
      });

      const allBalances = [
        ...store.wallet.balances,
        ...tokensNotInBalances.map((token) => Balance.create(token, "0")),
      ];

      return allBalances;
    });

    onMounted(async () => {
      await actions.refreshWalletBalances();
      await actions.refreshTokens();
    });

    return {
      balances,
    };
  },
  components: {},
});
</script>
