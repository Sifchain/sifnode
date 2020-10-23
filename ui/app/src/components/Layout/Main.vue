<template>
  <main class="main">
    <slot></slot>
  </main>
</template>

<script>
export default {

}
</script>

<style lang="scss">
.main {
  width: 100vw;
  min-height: calc(100vh - 81px); // account for menu height
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>