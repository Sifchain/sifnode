<template>
  <Modal class="with-wallet-container">
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
        :requestDialog="requestOpen"
        :connectCta="connectCta"
      ></slot>
    </template>
    <template v-slot:default>
      <div class="wallet-connect-container">
        <div class="top" />

        <div class="vstack">
          <EtheriumWalletPanel />
          <KeplrWalletPanel />
        </div>
      </div>
    </template>
  </Modal>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { useWalletButton } from "./useWalletButton";
import EtheriumWalletPanel from "./EtheriumWalletPanel.vue";
import KeplrWalletPanel from "./KeplrWalletPanel.vue";
import Modal from "@/components/shared/Modal.vue";

export default defineComponent({
  name: "WithWallet",
  components: { Modal, EtheriumWalletPanel, KeplrWalletPanel },
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
      connectCta,
    } = useWalletButton({
      connectType: props.connectType,
    });
    return {
      connected,
      connectedToEth,
      connectedToSif,
      connectCta,
    };
  },
});
</script>
<style lang="scss" scoped>
.wallet-connect-container {
  background: $c_gray_100;
  color: $c_gray_800;
}
.top {
  padding-top: 16px;
  padding-bottom: 16px;
}
.vstack {
  display: flex;
  flex-direction: column;
  & > * {
    margin: 1rem 1rem 0rem;
  }
  & > *:last-child {
    margin-bottom: 1rem;
  }
}
</style>
