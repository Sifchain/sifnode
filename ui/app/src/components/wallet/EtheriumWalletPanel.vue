<script lang="ts">
import { computed, defineComponent } from "vue";
import { useCore } from "@/hooks/useCore";
import SifButton from "@/components/shared/SifButton.vue";
import Icon from "@/components/shared/Icon.vue";

export default defineComponent({
  name: "EtheriumWalletController",
  components: {
    SifButton,
    Icon,
  },
  setup() {
    const { store, actions } = useCore();
    async function handleConnectClicked() {
      await actions.ethWallet.connectToWallet();
    }
    const address = computed(() => store.wallet.eth.address);
    const connected = computed(() => store.wallet.eth.isConnected);
    return {
      address,
      connected,
      handleConnectClicked,
    };
  },
});
</script>

<template>
  <div class="wrapper">
    <div v-if="connected">
      <p class="mb-2" v-if="address">{{ address }} <Icon icon="tick" /></p>
      <SifButton connect disabled>
        <img class="image" src="../../assets/metamask.png" />
        Metamask Connected
      </SifButton>
    </div>
    <SifButton connect v-else @click="handleConnectClicked"
      >Connect Metamask</SifButton
    >
  </div>
</template>

<style lang="scss" scoped>
.image {
  width: 23px;
  height: 100%;
  margin-top: 0px;
  margin-right: 12px;
}
</style>
