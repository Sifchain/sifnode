
<script lang="ts">
import ModalView from "./ModalView.vue";
import { ref, defineComponent } from "vue";

export default defineComponent({
  components: { ModalView },
  setup(props, context) {
    const isOpen = ref(false);

    return {
      isOpen,

      requestOpen() {
        isOpen.value = true;
        console.log("OPEN!");
      },
      requestClose(returnedData?: unknown) {
        console.log("requestClose");
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

