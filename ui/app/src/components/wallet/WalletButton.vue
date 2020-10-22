
<template>
  <button @click="handleClicked">
    <span v-if="!connected">Connect Wallet</span
    ><span v-else>{{ connectedText }}</span>
  </button>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useWalletButton } from "./useWalletButton";

export default defineComponent({
  name: "WalletButton",

  setup() {
    const { connected, connectedText, handleClicked } = useWalletButton({
      addrLen: 5,
    });
    return { connected, connectedText, handleClicked };
  },
});
</script>