<template>
  <div class="foo">
    <ConfirmationModal
      :requestClose="requestClose"
      :state="transactionState"
      :transactionHash="transactionHash"
      @confirmed="onConfirmed"
      :transactionStateMsg="transactionStateMsg"
      :confirmButtonText="'foobar'"
      :title="'You\'re beatiful'"
    >
      <template v-slot:common>
        <p class="text--normal">
          Unpegging <span class="text--bold">{{ amount }} {{ symbol }}</span>
        </p>
      </template>
    </ConfirmationModal>
    <button class="center" @click="transactionState = 'confirming'">go</button>
  </div>
</template>
<script lang="ts">
import { defineComponent, ref } from "vue";
import ConfirmationModal from "@/components/shared/ConfirmationModal.vue";
import { ConfirmState } from "../types";

export default defineComponent({
  components: {
    ConfirmationModal,
  },
  setup() {
    const pageState = {
      requestClose() {
        alert("close");
        pageState.transactionState.value = "selecting";
      },
      onConfirmed() {
        pageState.transactionState.value = "signing";
      },
      transactionState: ref<ConfirmState>("selecting"),
    };
    (window as any).pageState = pageState;
    return pageState;
  },
});
</script>
<style scoped>
.center {
  margin: 200px auto;
}
.foo {
  background: #666;
  width: 100%;
  height: 100vh; /* TODO: header height */
}
</style>