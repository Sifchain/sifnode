<template>
  <div class="loader" :class="classes">Loading...</div>
</template>
<script>
import { computed } from "@vue/reactivity";
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    success: { type: Boolean, default: false },
  },
  setup(props) {
    return {
      classes: computed(() => ({
        success: props.success,
      })),
    };
  },
});
</script>
<style lang="scss" scoped>
/* https://projects.lukehaas.me/css-loaders/ */
.loader,
.loader:after {
  border-radius: 50%;
  width: 60px;
  height: 60px;
}
.loader {
  font-size: 10px;
  position: relative;
  text-indent: -9999em;
  border-top: 2px solid $c_gold;
  border-right: 2px solid $c_gold;
  border-bottom: 2px solid $c_gold;
  border-left: 2px solid $c_gray_400;
  transform: translateZ(0);
  animation: load8 1.1s infinite linear;
  &.success {
    animation: none;
    border-top: 0;
    border-right: 0;
    border-bottom: 0;
    border-left: 0;
  }
}

@keyframes load8 {
  0% {
    -webkit-transform: rotate(0deg);
    transform: rotate(0deg);
  }
  100% {
    -webkit-transform: rotate(360deg);
    transform: rotate(360deg);
  }
}
</style>
