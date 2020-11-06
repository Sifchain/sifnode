<script lang="ts">
import { defineComponent, PropType } from "vue";
import Arrows from "./icons/Arrows.vue";
import CircleHalf from "./icons/CircleHalf.vue";
import PoolIcon from "./icons/PoolIcon.vue";
export default defineComponent({
  props: {
    icon: { type: String as PropType<"arrows" | "circle-half" | "pool"> },
  },
  components: { Arrows, CircleHalf, PoolIcon },
});
</script>
<template>
  <div>
    <Arrows v-if="icon === 'arrows'" />
    <CircleHalf v-if="icon === 'circle-half'" />
    <PoolIcon v-if="icon === 'pool'" />
  </div>
</template>