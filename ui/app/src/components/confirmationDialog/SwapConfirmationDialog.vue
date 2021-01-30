<template>
  <ConfirmationModal
    @confirmed="$emit('confirmswap')"
    :requestClose="requestClose"
    :isOpen="isOpen"
    :state="state"
  >
    <template v-slot:askbody>
      <div>THIS IS THE ASK BODY</div>
    </template>

    <template v-slot:signing>
      <div>signing</div>
    </template>

    <template v-slot:success>
      <div>success</div>
    </template>

    <template v-slot:error>
      <div>error</div>
    </template>
  </ConfirmationModal>
</template>

<script lang="ts">
import { PropType } from "vue";
import ConfirmationModal from "@/components/shared/ConfirmationModal.vue";
import { ConfirmState } from "../../types";

export default {
  components: {
    ConfirmationModal,
  },

  props: {
    isOpen: { type: Boolean, default: false },
    requestClose: Function,

    state: { type: String as PropType<ConfirmState>, default: "confirming" },
  },
};
</script>

<style>
</style>