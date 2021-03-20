<script lang="ts">
import Panel from "@/components/shared/Panel.vue";
import { defineComponent } from "vue";

export default defineComponent({
  components: { Panel },
  props: { isOpen: { type: Boolean, default: false }, requestClose: Function },
  emits: ["close"],
});
</script>

<template>
  <div class="modal">
    <teleport to="#portal-target">
      <transition name="fadein">
        <div class="backdrop" v-if="isOpen" @click="requestClose">
          <Panel class="modal-panel" v-if="isOpen" @click.stop>
            <slot :requestClose="requestClose"></slot>
            <div class="close" @click="requestClose">&times;</div>
          </Panel>
        </div>
      </transition>
    </teleport>
  </div>
</template>

<style lang="scss" scoped>
.activator {
  cursor: pointer;
}
.close {
  position: absolute;
  top: 16px;
  right: 16px;
  font-size: 32px;
  font-weight: normal;
  color: $c_text;
  cursor: pointer;
}
.modal-panel {
  position: relative;
  opacity: 1;
}

.fadein-leave-active {
  transition: opacity 0.3s ease-in-out 0.1s;
  & .modal-panel {
    transition: opacity 0.1s ease-in-out;
  }
}

.fadein-enter-active {
  transition: opacity 0.1s ease-in-out, transform 1s;
  & .modal-panel {
    transition: opacity 0.3s ease-in-out 0.1s;
  }
}

.fadein-enter-from,
.fadein-leave-to {
  opacity: 0;
  & .modal-panel {
    opacity: 0;
  }
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
  padding-top: 4rem; /* take into account header for centering */
}
</style>
