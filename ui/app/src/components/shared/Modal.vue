<template>
  <div class="modal">
    <div class="activator">
      <slot name="activator" :requestOpen="requestOpen"></slot>
    </div>
    <teleport to="body">
      <div class="backdrop" v-if="isOpen" @click="requestClose">
        <Panel :class="{ open: isOpen }" @click.stop>
          <div class="close" @click="requestClose">&times;</div>
          <slot :requestClose="requestClose"></slot>
        </Panel>
      </div>
    </teleport>
  </div>
</template>

<script lang="ts">
import Panel from "@/components/shared/Panel.vue";
import { ref, defineComponent } from "vue";

export default defineComponent({
  components: { Panel },
  setup(props, context) {
    const isOpen = ref(false);
    const selected = ref(null);
    return {
      isOpen,
      selected,
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
<style scoped>
.activator {
  cursor: pointer;
}
.close {
  cursor: pointer;
  background: none;
  border: none;
  position: absolute;
  top: 1rem;
  right: 1rem;
}
.backdrop {
  background: rgba(0, 0, 0, 0.1);
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
}
</style>

