<template>
  <div class="pool-list">
    <slot></slot>
  </div>
</template>

<script>
export default {};
</script>

<style scoped lang="scss">
.pool-list {
  height: 175px;
  overflow-y: auto;
  background: $c_white;
  border-radius: $br_sm;
}
</style>
