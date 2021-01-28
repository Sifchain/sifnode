<template>
  <ModalView
    :requestClose="requestClose"
    :isOpen="isOpen"
  > 
    <div class="modal-inner">
      <ConfirmationModalAsk 
        v-if="state === 'confirming'"
        @confirmed="$emit('confirmed')"
        title="Asking Confirmation"
      > 
        <template v-slot:body>
          <slot name="askbody"></slot>
        </template>
      </ConfirmationModalAsk>


      <ConfirmationModalSigning
        v-else
        :confirmed="confirmed"
        :failed="failed"
        :state="state"
        @closerequested="requestClose"
      >
        <template v-slot:signing>
          <slot name="signing"></slot>
        </template>

        <template v-slot:success>
          <slot name="success"></slot>
        </template>

        <template v-slot:error>
          <slot name="error"></slot>
        </template>
      </ConfirmationModalSigning>

    </div>


  <!-- <ConfirmationDialog
      @confirmswap="handleAskConfirmClicked"
      :transactionHash="transactionHash"
      :state="transactionState"
      :requestClose="requestTransactionModalClose"
      :priceMessage="priceMessage"
      :fromToken="fromSymbol"
      :fromAmount="fromAmount"
      :toAmount="toAmount"
      :toToken="toSymbol"
  /> -->
  </ModalView>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { computed } from "@vue/reactivity";

import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationModalAsk from "@/components/shared/ConfirmationModalAsk.vue";
import ConfirmationModalSigning from "@/components/shared/ConfirmationModalSigning.vue";

export type ConfirmState =
  | "selecting"
  | "confirming"
  | "signing"
  | "confirmed"
  | "failed";

export default defineComponent({
  props: { 
    isOpen: { type: Boolean, default: false }, 
    requestClose: Function,

    state: { type: String as PropType<ConfirmState>, default: "confirming" },
  },

  components: {
    ModalView,
    ConfirmationModalAsk,
    ConfirmationModalSigning,

    // ConfirmationDialog,
  },

  setup(props) {
    const confirmed = computed(() => {
      return props.state === "confirmed";
    });

    const failed = computed(() => {
      return props.state === "failed";
    });

    return {
      confirmed,
      failed
    };
  },
})
</script>

<style scoped lang="scss">
  .modal-inner {
    display: flex;
    flex-direction: column;
    padding: 30px 20px 20px 20px;
    min-height: 50vh;
  }
</style>