<template>
  <div class="home">
    <button v-if="!store.wallet.isConnected" @click="handleConnectClicked">
      Connect Wallet
    </button>

    <div v-if="store.wallet.isConnected">
      <p>Wallet Connected!</p>
      <table>
        <tr
          v-for="assetAmount in store.wallet.balances"
          :key="assetAmount.asset.symbol"
        >
          <td align="right">{{ assetAmount.toFixed(6) }}</td>
          <td align="left">{{ assetAmount.asset.symbol }}</td>
        </tr>
      </table>
      <button v-if="store.wallet.isConnected" @click="handleDisconnectClicked">
        Disconnect
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { Store, Actions } from "../../../core";
import { onMounted, inject } from "vue";

export default {
  name: "Wallet",
  setup() {
    const store = inject<Store>("store");
    const actions = inject<Actions>("actions");

    onMounted(async () => {
      await actions?.refreshWalletBalances();
    });

    async function handleConnectClicked() {
      await actions?.connectToWallet();
    }
    async function handleDisconnectClicked() {
      await actions?.disconnectWallet();
    }

    return { store, handleDisconnectClicked, handleConnectClicked };
  },
};
</script>
<style scoped>
table {
  margin: 0 auto;
}
</style>