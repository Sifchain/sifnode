<template>
  <div v-show="isActive">
    <slot></slot>
  </div>
</template>
<script>
import { inject } from "vue";
import { computed } from "@vue/reactivity";
export default {
  props: {
    title: String,
    slug: String,
  },
  setup(props) {
    const selectedTitle = inject("Tabs_selectedTitle");

    return {
      isActive: computed(() => props.title === selectedTitle.value),
    };
  },
};
</script>
