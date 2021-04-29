<template>
  <Icon class="arrow" :class="{ enabled, left }" icon="down" />
</template>
x

<script lang="ts">
import Icon from "@/components/shared/Icon.vue";
import { defineComponent } from "vue";

export default defineComponent({
  components: { Icon },
  props: {
    enabled: Boolean,
    left: Boolean,
  },
});
</script>
<style lang="scss" scoped>
.arrow {
  stroke: #c1c1c1;
  display: block;
  cursor: pointer;
  text-align: center !important;
  &.enabled {
    stroke: $c_gold;
  }
  &.left {
    text-align: left;
  }
}
</style>
