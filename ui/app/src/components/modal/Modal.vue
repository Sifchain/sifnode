<template>
  <div
    class="modal-backdrop"
    v-show="isOpen"
    :class="{ open: isOpen }"
    @click="$emit('close')"
  >
    <Panel :class="{ open: isOpen }" @click.stop>
      <button class="modal-dialog-close" @click="$emit('close')">x</button>
      <div class="modal-title" v-if="title">{{ title }}</div>
      <div class="modal-body">
        <slot />
      </div>
    </Panel>
  </div>
</template>

<script>
import { defineComponent } from "vue";
import Panel from "@/components/panel/Panel";

// https://medium.com/js-dojo/vue-js-manage-your-modal-window-s-effortlessly-using-eventbus-518977195eed
export default defineComponent({
  props: {
    isOpen: Boolean,
    title: String,
  },
  components: { Panel },
  setup() {
    return {};
  },
});
</script>

<style scoped>
.modal-dialog-close {
  background: none;
  border: none;
  position: absolute;
  top: 1rem;
  right: 1rem;
}
.modal-backdrop {
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

.modal-title {
  font-weight: bold;
  font-size: 1.3rem;
  margin-bottom: 1rem;
  color: rgb(0, 0, 0);
}
.modal-body {
  color: rgb(0, 0, 0);
}
</style>