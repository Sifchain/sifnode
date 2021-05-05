<template>
  <div class="price-calculation-cell" :data-handle="dataHandle">
    <slot></slot>
  </div>
</template>
<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: ["dataHandle"],
});
</script>
<style lang="scss" scoped>
.price-calculation-cell {
  width: 33%;
  flex: 1;
  text-align: center;
}
</style>
