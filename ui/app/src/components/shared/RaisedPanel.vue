<template>
  <div class="raised-panel">
    <slot></slot>
  </div>
</template>
<style lang="scss" scoped>
.raised-panel {
  padding: 4px 15px 15px 15px;
  border-radius: $br_sm;
  border: 1px solid $c_gray_100;
  background: $g_gray_reverse;
  color: $c_gray_700;
}
</style>
