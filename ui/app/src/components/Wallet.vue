<template>
  <div class="home">
    <table>
      <tr
        v-for="assetAmount in state.availableAssetAccounts"
        :key="assetAmount.asset.symbol"
      >
        <td align="right">{{ assetAmount.toFixed(6) }}</td>
        <td align="left">{{ assetAmount.asset.symbol }}</td>
      </tr>
    </table>
  </div>
</template>

<script lang="ts">
import { State, UseCases } from "../../../core";
import { onMounted, inject } from "vue";

export default {
  name: "Wallet",
  setup() {
    const state = inject<State>("state");
    const usecases = inject<UseCases>("usecases");

    onMounted(async () => {
      if (usecases) {
        await usecases.updateListOfAvailableTokens();
      }
    });

    return { state };
  },
};
</script>
