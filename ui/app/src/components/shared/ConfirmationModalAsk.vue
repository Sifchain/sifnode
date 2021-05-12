<template>
  <div class="ask-wrapper">
    <h3 class="title mb-10" data-handle="confirmation-modal-title">
      {{ title }}
    </h3>
    <div class="ask-body">
      <slot name="body"></slot>
    </div>
    <SifButton block primary class="confirm-btn" @click="$emit('confirmed')">{{
      confirmButtonText
    }}</SifButton>
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
  },
});
</script>

<style lang="scss" scoped>
.ask-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.title {
  font-size: $fs_lg;
  color: $c_text;
  margin-bottom: 1rem;
  text-align: left;
  font-weight: 700;
  font-style: italic;
}
.ask-body {
  flex: 1;
}
.confirm-btn {
  margin-top: auto;
}
</style>
