<template>
  <div class="home">
    <button
      v-if="!store.wallet.etheriumIsConnected"
      @click="handleConnectClicked"
    >
      Connect Wallet
    </button>

    <div v-if="store.wallet.etheriumIsConnected">
      <p>Wallet Connected!</p>
      <table>
        <tr
          v-for="assetAmount in store.wallet.balances"
          :key="assetAmount.asset.symbol"
        >
          <td align="left">{{ assetAmount.asset.symbol }}</td>
          <td align="right">{{ assetAmount.toFixed(2) }}</td>
        </tr>
      </table>
      <button
        v-if="store.wallet.etheriumIsConnected"
        @click="handleDisconnectClicked"
      >
        Disconnect
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { onMounted } from "vue";
import { useCore } from "../hooks/useCore";

export default {
  name: "Wallet",
  setup() {
    const { store, actions } = useCore();

    onMounted(() => {
      actions.init();
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