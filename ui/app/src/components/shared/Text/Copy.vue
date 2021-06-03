<template>
  <div class="container">
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>
.container {
  font-size: 12px; /* make $variable */
  line-height: 1.5; /* make $variable */
  text-align: left;
  font-weight: 400; /* make $variable */
}
</style>
