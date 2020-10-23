<template>
  <div class="sif-app__wrapper">
    <div class="sif-app__before">
      <slot name="before"></slot>
    </div>

    <div class="sif-app">
      <div class="sif-app__header">
        <slot name="header"></slot>
      </div>

      <div class="sif-app__body">
        <slot></slot>
      </div>

      <div class="sif-app__footer">
        <slot name="footer"></slot>
      </div>
    </div> <!-- main app card -->

    <div class="sif-app__after">
      <slot name="after"></slot>
    </div>
  </div> <!-- wrapper -->
</template>

<script>
export default {

}
</script>

<style lang="scss">
.sif-app {
  width: 410px;
  border-radius: $br_lg;
  overflow: hidden;

  &__header {
    padding: 20px 20px 10px 20px;
    background: $c_gray_900;
    color: $c_gray_400;
  }

  &__body {
    padding: 15px;
    background: $g_gray;
    color: $c_gray_700;
  }
}
</style>