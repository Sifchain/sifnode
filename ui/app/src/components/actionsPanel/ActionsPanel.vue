
<script lang="ts">
import { defineComponent } from "vue";
import WithWallet from "@/components/wallet/WithWallet.vue";
import SifButton from "@/components/shared/SifButton.vue";
import Icon from "@/components/shared/Icon.vue";

export default defineComponent({
  components: {
    WithWallet,
    SifButton,
    Icon,
  },
  props: {
    nextStepAllowed: Boolean,
    nextStepMessage: String,
  },
  emits: ["nextstepclick"],
  setup(_, { emit }) {
    return {
      handleNextStepClicked() {
        emit("nextstepclick");
      },
    };
  },
});
</script>

<template>
  <div class="actions">
    <WithWallet>
      <template v-slot:disconnected="{ requestDialog }">
        <div class="wallet-status">
          No wallet connected <Icon icon="cross" />
        </div>
        <SifButton primary block @click="requestDialog">
          Connect Wallet
        </SifButton>
      </template>
      <template v-slot:connected="{ connectedText }"
        ><div>
          <div class="wallet-status">
            Connected to {{ connectedText }} <Icon icon="tick" />
          </div>
          <SifButton
            v-if="nextStepMessage"
            block
            primary
            :disabled="!nextStepAllowed"
            @click="handleNextStepClicked"
          >
            {{ nextStepMessage }}
          </SifButton>
        </div></template
      >
    </WithWallet>
  </div>
</template>

<style lang="scss" scoped>
.actions {
  padding-top: 1rem;
}
.wallet-status {
  font-size: $fs_sm;
}
</style>
