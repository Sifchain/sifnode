<script lang="ts">
import { defineComponent, PropType } from "vue";
import Arrows from "./icons/Arrows.vue";
import CircleHalf from "./icons/CircleHalf.vue";
export default defineComponent({
  props: {
    icon: { type: String as PropType<"arrows" | "circle-half"> },
  },
  components: { Arrows, CircleHalf },
});
</script>
<template>
  <div>
    <Arrows v-if="icon === 'arrows'" />
    <CircleHalf v-if="icon === 'circle-half'" />
  </div>
</template>