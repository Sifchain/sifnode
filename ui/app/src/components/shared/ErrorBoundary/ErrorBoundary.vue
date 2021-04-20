<template>
  <div class="guard-container">
    <slot></slot>
    <div v-if="error" class="guard">
      <slot name="fallback"></slot>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "@vue/runtime-core";

export default defineComponent({
  props: ["error"],
});
</script>

<style scoped lang="scss">
.guard-container {
}
.guard {
  position: absolute;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  background: pink;
}
</style>
