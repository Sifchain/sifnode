<template>
  <div v-if="!store.wallet.eth.isConnected">
    <div v-if="!connecting">
      <p>Choose the wallet you want to connect to the etherium blockchain</p>
      <button @click="handleMetamaskClicked">Connect with metamask</button>
    </div>
    <div v-else>Connecting...</div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { ref } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { ModalBus } from "@/components/modal/ModalBus";
import WalletConnectedVue from "./WalletConnected.vue";
export default defineComponent({
  name: "WalletConnectSelector",
  setup() {
    const { store, actions } = useCore();
    const connecting = ref<boolean>(false);

    async function handleMetamaskClicked() {
      connecting.value = true;

      await actions.connectToWallet();

      connecting.value = false;

      ModalBus.emit("open", {
        title: "Your Wallet",
        component: WalletConnectedVue,
      });
    }

    return {
      store,
      handleMetamaskClicked,
      connecting,
    };
  },
});
</script>