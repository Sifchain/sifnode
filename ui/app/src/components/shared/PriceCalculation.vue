<template>
  <div class="price-calculation">
    <slot></slot>
  </div>
</template>

<script>
export default {
}
</script>

<style lang="scss">
.price-calculation {
  height: 120px;
  margin-bottom: 10px;
  border-radius: $br_sm;
  background: $c_white;
  padding: 16px;
  overflow-y: auto;
}
</style>