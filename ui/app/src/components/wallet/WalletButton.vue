<template>
  <button @click="handleClicked">
    <span v-if="!connected">Connect Wallet</span
    ><span v-else>{{ connectedText }}</span>
  </button>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { computed } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { ModalBus } from "@/components/modal/ModalBus";
import WalletConnectSelectorVue from "./WalletConnectSelector.vue";
import WalletConnectedVue from "./WalletConnected.vue";
function shorten(str: string) {
  return str.slice(0, 5) + "...";
}
export default defineComponent({
  name: "WalletButton",

  setup() {
    const { store } = useCore();

    async function handleClicked() {
      const dialog = !store.wallet.eth.isConnected
        ? {
            component: WalletConnectSelectorVue,
            title: "Select your wallet",
          }
        : {
            component: WalletConnectedVue,
            title: "Your wallet",
          };

      ModalBus.emit("open", dialog);
    }

    const connected = computed(() => store.wallet.eth.isConnected);
    const connectedText = computed(() =>
      [store.wallet.eth.address, store.wallet.sif.address]
        .filter(Boolean)
        .map(shorten)
        .join(", ")
    );

    return { connected, connectedText, handleClicked };
  },
});
</script>