<template>
  <div class="home">
    <p>Hello World</p>
    <div v-for="balance in balances" :key="balance">
      {{ balance }}
    </div>
  </div>
</template>

<script lang="ts">
import Vue, { VueConstructor } from "vue";
import { api } from "sif-core";
import { AssetAmount } from "sif-core/dist/entities";

export default (Vue as VueConstructor<Vue & { balances: string[] }>).extend({
  name: "Home",
  data() {
    return { balances: [] } as { balances: string[] };
  },
  async mounted() {
    const balances = await api.walletService.getAssetBalances();
    this.balances = balances.map((amount: AssetAmount) => {
      const str = amount.asset.symbol + ":" + amount.toFixed();
      return str;
    });
  },
});
</script>
