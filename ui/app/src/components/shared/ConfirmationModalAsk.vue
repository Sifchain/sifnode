<template>
  <div class="modal-inner">
    <div class="ask-wrapper">
      <h3 class="title mb-10">{{ title }}</h3>
      <div class="ask-body">
        <slot name="body"></slot>
      </div>
      <SifButton
        block
        primary
        class="confirm-btn"
        @click="$emit('confirmed')"
        >{{ confirmButtonText }}</SifButton
      >
    </div>
  </div>
</template>

<script lang="ts">
import SifButton from "@/components/shared/SifButton.vue";
import { defineComponent } from "vue";

export default defineComponent({
  inheritAttrs: false,
  components: {
    SifButton,
  },

  props: {
    title: String,
    confirmButtonText: String,
    onConfirmed: Function,
  },
});
</script>

<style lang="scss" scoped>
.modal-inner {
  display: flex;
  flex-direction: column;
  padding: 30px 20px 20px 20px;
  min-height: 50vh;
}
.ask-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.title {
  font-size: $fs_lg;
  color: $c_text;
  margin-bottom: 2rem;
  text-align: left;
  font-weight: 700;
}
.ask-body {
  flex: 1;
}
.confirm-btn {
  margin-top: auto;
}
</style>