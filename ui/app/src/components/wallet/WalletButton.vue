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
import WalletConnectDialog from "./WalletConnectDialog.vue";

function shorten(str: string) {
  return str.slice(0, 5) + "...";
}
export default defineComponent({
  name: "WalletButton",

  setup() {
    const { store } = useCore();

    async function handleClicked() {
      ModalBus.emit("open", { component: WalletConnectDialog });
    }

    const connected = computed(
      () => store.wallet.eth.isConnected || store.wallet.sif.isConnected
    );
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