<script lang="ts">
import { computed, defineComponent } from "vue";
import { useCore } from "@/hooks/useCore";
import BalanceTable from "./BalanceTable.vue";
import SifButton from "@/components/shared/SifButton.vue";
import Icon from "@/components/shared/Icon.vue";

function useKeplrWallet() {
  const { store, actions } = useCore();
  async function handleDisconnectClicked() {
    await actions.wallet.disconnectWallet();
  }
  async function handleConnectClicked() {
    try {
      await actions.wallet.connectToWallet();
    } catch (error) {
      console.log("ui", error);
    }
  }
  const address = computed(() => store.wallet.sif.address);
  const balances = computed(() => store.wallet.sif.balances);
  const connected = computed(() => store.wallet.sif.isConnected);
  return {
    address,
    balances,
    connected,
    handleConnectClicked,
    handleDisconnectClicked,
  };
}
export default defineComponent({
  name: "KeplrWalletController",
  components: { BalanceTable, SifButton, Icon },
  setup() {
    const {
      address,
      balances,
      connected,
      handleConnectClicked,
      handleDisconnectClicked,
    } = useKeplrWallet();
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
        >Disconnect Keplr</SifButton
      >
    </div>
    <SifButton connect v-else @click="handleConnectClicked">Keplr</SifButton>
  </div>
</template>
