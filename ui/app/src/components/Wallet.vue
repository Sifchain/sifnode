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
import { onMounted } from "vue";
import { useCore } from "../core/useCore";

export default {
  name: "Wallet",
  setup() {
    const { store, actions } = useCore();

    onMounted(async () => {
      await actions.refreshWalletBalances();
    });

    async function handleConnectClicked() {
      await actions.connectToWallet();
    }
    async function handleDisconnectClicked() {
      await actions.disconnectWallet();
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