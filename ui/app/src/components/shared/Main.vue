<template>
  <main class="main">
    <slot></slot>
  </main>
</template>

<script>
export default {};
</script>

<style lang="scss" scoped>
.main {
  width: 100vw;
  min-height: calc(100vh - 80px); // account for menu height
  margin-top: 80px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: url("../../assets/World_Background_opt.jpg");
  background-size: cover;
  background-position: bottom center;
}
</style>
