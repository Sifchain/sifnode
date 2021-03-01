<template>
  <div class="wrapper">
    <div class="before">
      <slot name="before"></slot>
    </div>

    <div class="sif-app">
      <div class="header">
        <slot name="header"></slot>
      </div>

      <div class="body">
        <slot></slot>
      </div>

      <div class="footer">
        <slot name="footer"></slot>
      </div>
    </div>
    <!-- main app card -->

    <div class="after">
      <slot name="after"></slot>
    </div>
  </div>
  <!-- wrapper -->
</template>

<script>
export default {};
</script>

<style lang="scss" scoped>
.sif-app {
  width: 410px;
  border-radius: $br_lg;
  overflow: hidden;

  .header {
    padding: 20px 20px 10px 20px;
    background: $c_gray_900;
    color: $c_gray_400;
  }

  .body {
    padding: 15px;
    background: $g_gray;
    color: $c_gray_700;
  }
}
</style>
