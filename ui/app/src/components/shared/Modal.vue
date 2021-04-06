<script lang="ts">
import ModalView from "./ModalView.vue";
import { ref, defineComponent } from "vue";

export default defineComponent({
  components: { ModalView },
  props: ["open"],
  setup(props, context) {
    const isOpen = ref(props.open || false);

    return {
      isOpen,

      requestOpen() {
        isOpen.value = true;
      },
      requestClose(returnedData?: unknown) {
        isOpen.value = false;
        context.emit("close", returnedData);
      },
    };
  },
  emits: ["close"],
});
</script>

<template>
  <slot name="activator" :requestOpen="requestOpen"></slot>
  <ModalView :isOpen="isOpen" :requestClose="requestClose"
    ><slot :requestClose="requestClose"></slot
  ></ModalView>
</template>
