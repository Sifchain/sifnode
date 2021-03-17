<template>
  <div class="raised-panel-column">
    <slot></slot>
  </div>
</template>
<style lang="scss" scoped>
.raised-panel-column {
  text-align: left;
  display: flex;
  flex-direction: column;

  flex-grow: 1;
}
</style>
