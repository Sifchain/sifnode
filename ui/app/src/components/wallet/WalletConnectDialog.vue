<template>
  <EtheriumWalletPanel />
  <hr />
  <SifWalletPanel />

  <hr />
  <button @click="handleClose">Close</button>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import EtheriumWalletPanel from "./EtheriumWalletPanel.vue";
import SifWalletPanel from "./SifWalletPanel.vue";

export default defineComponent({
  name: "WalletConnectDialog",
  emits: ["close"],
  components: { EtheriumWalletPanel, SifWalletPanel },
  setup(_, { emit }) {
    function handleClose() {
      emit("close");
    }

    return {
      handleClose,
    };
  },
});
</script>
<style scoped>
.h-stack {
  display: flex;
}
</style>