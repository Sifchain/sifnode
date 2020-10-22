<template>
  <div
    class="modal-backdrop"
    v-show="isOpen"
    :class="{ open: isOpen }"
    @click="$emit('close')"
  >
    <div class="modal-dialog" :class="{ open: isOpen }" @click.stop>
      <button class="modal-dialog-close" @click="$emit('close')">x</button>
      <div class="modal-title" v-if="title">{{ title }}</div>
      <div class="modal-body">
        <slot />
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent } from "vue";

// https://medium.com/js-dojo/vue-js-manage-your-modal-window-s-effortlessly-using-eventbus-518977195eed
export default defineComponent({
  props: {
    isOpen: Boolean,
    title: String,
  },
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
  top: 0;
  right: 0;
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
.modal-dialog {
  position: relative;
  width: 30rem;
  background: rgb(255, 255, 255);
  padding: 1.5rem 2rem;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1),
    0 4px 6px -2px rgba(0, 0, 0, 0.05);
  border-radius: 0.3rem;
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