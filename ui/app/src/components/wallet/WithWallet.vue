
<template>
  <Modal>
    <template v-slot:activator="{ requestOpen }">
      <slot
        v-if="!connected"
        name="disconnected"
        :connectedToEth="connectedToEth"
        :connectedToSif="connectedToSif"
        :requestDialog="requestOpen"
        :connectCta="connectCta"
      ></slot>
      <slot
        v-if="connected"
        name="connected"
        :connectedToEth="connectedToEth"
        :connectedToSif="connectedToSif"
        :connectedText="connectedText"
        :requestDialog="requestOpen"
        :connectCta="connectCta"
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
  props: {
    connectType: String as PropType<
      "connectToAny" | "connectToAll" | "connectToSif"
    >,
  },
  setup(props) {
    const {
      connected,
      connectedToEth,
      connectedToSif,
      connectedText,
      connectCta,
    } = useWalletButton({
      addrLen: 10,
      connectType: props.connectType,
    });
    return {
      connected,
      connectedToEth,
      connectedToSif,
      connectedText,
      connectCta,
    };
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