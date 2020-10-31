<template>
  <div class="modal">
    <slot name="activator" :requestOpen="requestOpen"></slot>
    <teleport to="#portal-target">
      <transition name="foo">
        <div class="backdrop" v-if="isOpen" @click="requestClose">
          <Panel class="panel" @click.stop>
            <div class="close" @click="requestClose">&times;</div>
            <slot :requestClose="requestClose"></slot>
          </Panel>
        </div>
      </transition>
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
<style lang="scss" scoped>
.activator {
  cursor: pointer;
}
.close {
  position: absolute;
  top: 16px;
  right: 20px;
  font-size: 32px;
  font-weight: normal;
  color: $c_text;
  cursor: pointer;
}

.foo-leave-active {
  transition: opacity 0.2s ease-in-out;
}

.foo-enter-active {
  transition: opacity 0.1s ease-in-out;
  & .panel {
    transition: transform 0.05s ease-in-out;
  }
}

.foo-enter-from {
  opacity: 0;
  & .panel {
    transform: translateY(20px);
  }
}

.foo-leave-to {
  opacity: 0;
}

.backdrop {
  background: rgba(0, 0, 0, 0.4);
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

