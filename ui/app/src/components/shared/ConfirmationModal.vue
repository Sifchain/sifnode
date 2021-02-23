<template>
  <ModalView :requestClose="requestClose" :isOpen="isOpen">
    <div class="modal-inner">
      <ConfirmationModalAsk
        v-if="state === 'confirming'"
        @confirmed="$emit('confirmed')"
        :title="title"
        :confirmButtonText="confirmButtonText"
      >
        <template v-slot:body>
          <slot name="selecting"></slot>
        </template>
      </ConfirmationModalAsk>

      <ConfirmationModalSigning
        v-else
        :state="state"
        @closerequested="requestClose"
        :transactionHash="transactionHash"
        :transactionStateMsg="transactionStateMsg"
      >
        <template v-slot:signing>
          <slot :name="!!this.$slots.signing ? 'signing' : 'common'"></slot>
        </template>

        <template v-slot:confirmed>
          <slot :name="!!this.$slots.signing ? 'confirmed' : 'common'"></slot>
        </template>

        <template v-slot:rejected>
          <slot :name="!!this.$slots.signing ? 'rejected' : 'common'"></slot>
        </template>

        <template v-slot:failed>
          <slot :name="!!this.$slots.signing ? 'failed' : 'common'"></slot>
        </template>
      </ConfirmationModalSigning>
    </div>
  </ModalView>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { computed } from "@vue/reactivity";
import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationModalAsk from "@/components/shared/ConfirmationModalAsk.vue";
import ConfirmationModalSigning from "@/components/shared/ConfirmationModalSigning.vue";
import { ConfirmState } from "../../types";

export default defineComponent({
  inheritAttrs: false,
  props: {
    // Function to request the window is closed this function must reset the confirmation state to selecting
    requestClose: Function,

    // Confirmation state: "selecting" | "confirming" | "signing" | "confirmed" | "rejected" | "failed";
    // This component acts on this state to determine which panel to show
    state: { type: String as PropType<ConfirmState>, default: "confirming" },

    // The text on the 'confirm' button
    confirmButtonText: String,

    // The title of the ask modal
    title: String,

    // TODO: Revisit if we need these here or can make them part of the content
    transactionHash: String,
    transactionStateMsg: String,
  },

  setup(props) {
    const isOpen = computed(() => {
      return [
        "approving",
        "confirming",
        "signing",
        "failed",
        "rejected",
        "confirmed",
      ].includes(props.state);
    });

    return {
      isOpen,
    };
  },

  components: {
    ModalView,
    ConfirmationModalAsk,
    ConfirmationModalSigning,
  },
});
</script>

<style scoped lang="scss">
.modal-inner {
  display: flex;
  flex-direction: column;
  padding: 30px 20px 20px 20px;
  min-height: 50vh;
}
</style>