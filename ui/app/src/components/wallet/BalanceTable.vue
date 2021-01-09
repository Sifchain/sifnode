<template>
  <table>
    <tr v-for="balance in balanceList" :key="balance.symbol">
      <td align="left">{{ balance.symbol }}</td>
      <td align="right">{{ balance.amount }}</td>
    </tr>
  </table>
</template>

<script lang="ts">
import { computed } from "@vue/reactivity";
import { AssetAmount } from "ui-core";
import { defineComponent } from "vue";
import { labelDecorator } from "@/utils/labelDecorator";

export default defineComponent({
  props: ["balances"],
  setup(props) {
    return {
      balanceList: computed(() =>
        props.balances.map((balance: AssetAmount) => {
          return {
            symbol: labelDecorator(balance.asset.symbol),
            amount: balance.toFixed(),
          };
        })
      ),
    };
  },
});
</script>
<style scoped>
table {
  width: 100%;
}
</style>