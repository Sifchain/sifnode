
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
import { defineComponent, PropType } from "vue";
import { useWalletButton } from "./useWalletButton";
import EtheriumWalletPanel from "./EtheriumWalletPanel.vue";
import SifWalletPanel from "./SifWalletPanel.vue";

import Modal from "@/components/shared/Modal.vue";
export default defineComponent({
  name: "WithWallet",
  components: { Modal, EtheriumWalletPanel, SifWalletPanel },
  setup() {
    const { connected, connectedText } = useWalletButton({
      addrLen: 10,
    });
    return { connected, connectedText };
  },
});
</script>
<style lang="scss" scoped>
.vstack {
  display: flex;
  flex-direction: column;
  padding-top: 1rem;
  & > * {
    margin: 2rem 2rem 1rem;
  }
  & > *:last-child {
    margin-bottom: 2rem;
  }
}
</style>