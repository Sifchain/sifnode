<template>
  <div class="input-group">
    <slot></slot>
  </div>
</template>

<script>
export default {};
</script>

<style lang="scss" scoped>
.input-group {
  margin-bottom: 20px;
  padding: 4px 15px 15px 15px;
  border-radius: $br_sm;
  background: $g_gray_reverse;
  color: $c_gray_700;
}
</style>
