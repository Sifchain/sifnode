<template>
  <table>
    <tr v-for="assetAmount in balances" :key="assetAmount.asset.symbol">
      <td align="left">{{ assetAmount.asset.symbol }}</td>
      <td align="right">{{ assetAmount.toFixed() }}</td>
    </tr>
  </table>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: ["balances"],
});
</script>
<style scoped>
table {
  width: 100%;
}
</style>