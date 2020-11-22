<script lang="ts">
import { defineComponent, PropType } from "vue";
import CircleArrows from "./icons/CircleArrows.vue";
import CircleHalf from "./icons/CircleHalf.vue";
import PoolIcon from "./icons/PoolIcon.vue";
import DownArrow from "./icons/DownArrow.vue";
import Tick from "./icons/Tick.vue";
import Cross from "./icons/Cross.vue";

export default defineComponent({
  props: {
    icon: {
      type: String as PropType<
        "arrows" | "circle-half" | "pool" | "down" | "tick" | "cross"
      >,
    },
  },
  components: { Cross, Tick, DownArrow, CircleArrows, CircleHalf, PoolIcon },
});
</script>

<template>
  <span>
    <CircleArrows v-if="icon === 'circle-arrows'" />
    <CircleHalf v-if="icon === 'circle-half'" />
    <PoolIcon v-if="icon === 'pool'" />
    <DownArrow v-if="icon === 'down'" />
    <Tick v-if="icon === 'tick'" />
    <Cross v-if="icon === 'cross'" />
  </span>
</template>