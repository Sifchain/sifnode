<template>
  <Icon class="arrow" :class="{ enabled }" icon="down" />
</template>

<script lang="ts">
import Icon from "@/components/shared/Icon.vue";
import { defineComponent } from "vue";

export default defineComponent({
  components: { Icon },
  props: {
    enabled: Boolean,
  },
});
</script>
<style lang="scss" scoped>
.arrow {
  stroke: #c1c1c1;
  &.enabled {
    stroke: $c_gold;
  }
}
</style>