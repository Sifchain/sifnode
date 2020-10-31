
<template>
  <Modal>
    <template v-slot:activator="{ requestOpen }">
      <slot
        v-if="!connected"
        name="disconnected"
        :requestDialog="requestOpen"
      ></slot>
      <slot
        v-if="connected"
        name="connected"
        :connectedText="connectedText"
        :requestDialog="requestOpen"
      ></slot>
    </template>
    <template v-slot:default>
      <div class="vstack">
        <EtheriumWalletPanel />
        <SifWalletPanel />
      </div>
    </template>
  </Modal>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useWalletButton } from "./useWalletButton";
import EtheriumWalletPanel from "./EtheriumWalletPanel.vue";
import SifWalletPanel from "./SifWalletPanel.vue";

import Modal from "@/components/shared/Modal.vue";
export default defineComponent({
  name: "WithWallet",
  components: { Modal, EtheriumWalletPanel, SifWalletPanel },
  setup() {
    const { connected, connectedText } = useWalletButton({
      addrLen: 5,
    });
    return { connected, connectedText };
  },
});
</script>
<style lang="scss" scoped>
.vstack {
  display: flex;
  flex-direction: column;
  & > * {
    margin-bottom: 1rem;
  }
}
</style>