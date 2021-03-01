<template>
  <div class="swap-loader" :class="{ loading: loading, success: success }">
    <div class="ring"></div>
    <slot></slot>
  </div>
</template>

<script>
export default {
  props: {
    loading: Boolean,
    success: Boolean,
  },
};
</script>

<style scoped lang="scss">
.swap-loader {
  position: relative;
  width: 20px;
  height: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-style: normal !important;
  background: $c_gray_400;
  color: $c_white;
  border-radius: 20px;

  &.success {
    background: $c_success;
  }

  &.loading {
    background: $c_gold;

    .ring {
      position: absolute;
      width: 36px;
      height: 36px;
      left: -8px;
      top: -8px;
      border: 2px solid $c_gold;
      border-radius: 40px;
      border-top-color: transparent;
      animation: 1s rotate infinite linear;
    }
  }
}

@keyframes rotate {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
