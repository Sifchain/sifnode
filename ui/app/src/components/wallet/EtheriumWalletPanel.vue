<script lang="ts">
import { computed, defineComponent } from "vue";
import { useCore } from "@/hooks/useCore";
import BalanceTable from "./BalanceTable.vue";

function useEtheriumWallet() {
  const { store, actions } = useCore();

  async function handleDisconnectClicked() {
    await actions.ethWallet.disconnectWallet();
  }

  async function handleConnectClicked() {
    await actions.ethWallet.connectToWallet();
  }

  const address = computed(() => store.wallet.eth.address);
  const balances = computed(() => store.wallet.eth.balances);
  const connected = computed(() => store.wallet.eth.isConnected);

  return {
    address,
    balances,
    connected,
    handleConnectClicked,
    handleDisconnectClicked,
  };
}

export default defineComponent({
  name: "EtheriumWalletController",
  components: { BalanceTable },
  setup() {
    const {
      address,
      balances,
      connected,
      handleConnectClicked,
      handleDisconnectClicked,
    } = useEtheriumWallet();
    return {
      address,
      balances,
      connected,
      handleConnectClicked,
      handleDisconnectClicked,
    };
  },
});
</script>

<template>
  <div>
    <div v-if="connected">
      <p>{{ address }}</p>
      <BalanceTable :balances="balances" />
      <button @click="handleDisconnectClicked">DisconnectWallet</button>
    </div>
    <button v-else @click="handleConnectClicked">Connect to Metamask</button>
  </div>
</template>