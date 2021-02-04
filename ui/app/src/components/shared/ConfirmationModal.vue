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
        <template v-slot:approving>
          <slot :name="!!this.$slots.signing ? 'confirmed' : 'common'"></slot>
        </template>

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
import { TransactionStatus } from "ui-core";

export type ConfirmState =
  // selecting values?
  | "selecting"
  | "approving"
  | "confirming"
  | "signing"
  | "confirmed"
  | "rejected"
  | "failed";

export default defineComponent({
  inheritAttrs: false,
  props: {
    // Function to request the window is closed this function must reset the confirmation state to selecting
    requestClose: Function,

    // TxStatus
    txStatus: {
      type: Object as PropType<TransactionStatus>,
      default: "confirming",
    },

    // The text on the 'confirm' button
    confirmButtonText: String,

    // The title of the ask modal
    title: String,

    // is the dialog open
    isOpen: Boolean,

    onConfirmed: Function,
  },

  setup(props) {
    //   const isOpen = computed(() => {
    //     return [
    //       "confirming",
    //       "signing",
    //       "failed",
    //       "rejected",
    //       "confirmed",
    //     ].includes(props.state);
    //   });

    //   return {
    //     isOpen,
    //   };
    return {
      state: computed<ConfirmState>(() => {
        if (!props.txStatus) {
          return "selecting";
        }

        if (props.txStatus.state === "failed") return "failed";
        return "failed";
      }),
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