<template>
  <ModalView
    :requestClose="requestClose"
    :isOpen="isOpen"
  > 
    <div class="modal-inner">

      <ConfirmationModalAsk 
        v-if="state === 'confirming'"
        @confirmed="$emit('confirmed')"
        :title="title"
        :confirmButtonText="confirmButtonText"
      > 
        <template v-slot:body>
          <slot name="askbody"></slot>
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
          <slot name="signing"></slot>
        </template>

        <template v-slot:confirmed>
          <slot name="confirmed"></slot>
        </template>

        <template v-slot:rejected>
          <slot name="rejected"></slot>
        </template>

        <template v-slot:failed>
          <slot name="failed"></slot>
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

export type ConfirmState =
  | "selecting"
  | "confirming"
  | "signing"
  | "confirmed"
  | "rejected"
  | "failed";

export default defineComponent({
  props: { 
    isOpen: { type: Boolean, default: false }, 
    requestClose: Function,
    state: { type: String as PropType<ConfirmState>, default: "confirming" },
    confirmButtonText: String,
    title: String,
    transactionHash: String,
    transactionStateMsg: String,
  },

  components: {
    ModalView,
    ConfirmationModalAsk,
    ConfirmationModalSigning,
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