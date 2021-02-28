<script lang="ts">
import { computed, defineComponent } from "vue";
import { useCore } from "@/hooks/useCore";
import SifButton from "@/components/shared/SifButton.vue";
import Icon from "@/components/shared/Icon.vue";

export default defineComponent({
  name: "KeplrWalletController",
  components: { SifButton, Icon },
  setup() {
    const { store, actions } = useCore();
    function formatAddress(address: string) {
      return !address || address.length < 4
        ? ""
        : address.substring(0, 7) +
            "..." +
            address.substring(address.length - 4);
    }
    async function handleConnectClicked() {
      try {
        await actions.wallet.connectToWallet();
      } catch (error) {
        console.log("KeplrWalletController", error);
      }
    }
    const address = computed(() => store.wallet.sif.address);
    const connected = computed(() => store.wallet.sif.isConnected);
    return {
      address,
      connected,
      formatAddress,
      handleConnectClicked,
    };
  },
});
</script>

<template>
  <div class="wrapper">
    <div v-if="connected">
      <img class="image" src="../../assets/keplr.jpg" />
      <p class="mb-2" v-if="address">
        {{ formatAddress(address) }} <Icon icon="tick" />
      </p>
    </div>
    <SifButton connect v-else @click="handleConnectClicked">Keplr</SifButton>
  </div>
</template>

<style lang="scss" scoped>
.image {
  height: 32px;
}
</style>
