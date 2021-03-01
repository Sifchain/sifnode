<script lang="ts">
import { computed, defineComponent } from "vue";
import { useCore } from "@/hooks/useCore";
import BalanceTable from "./BalanceTable.vue";
import SifButton from "@/components/shared/SifButton.vue";
import Icon from "@/components/shared/Icon.vue";

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
  components: {
    // BalanceTable,
    SifButton,
    Icon,
  },
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
  <div class="wrapper">
    <div v-if="connected">
      <p class="mb-2" v-if="address">{{ address }} <Icon icon="tick" /></p>
      <!-- <BalanceTable :balances="balances" /> -->
      <SifButton connect active @click="handleDisconnectClicked"
        >Disconnect Metamask</SifButton
      >
    </div>
    <SifButton connect v-else @click="handleConnectClicked">Metamask</SifButton>
  </div>
</template>
