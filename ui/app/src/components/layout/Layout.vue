<template>
  <div class="layout">
    <slot></slot>
  </div>
</template>

<style scoped>
.layout {
  background: #f9f9f9;
  box-sizing: border-box;
  padding-top: 3rem;
  width: 100%;
  height: calc(100vh - 4rem); /* TODO: header height */
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: flex-start;
}
</style>