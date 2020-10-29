<template>
    <div class="panel">
      <div class="header">
        <slot name="header"></slot>
      </div>

      <div class="body">
        <slot></slot>
      </div>

      <div class="footer">
        <slot name="footer"></slot>
      </div>
    </div> <!-- main app card -->
</template>

<style lang="scss" scoped>
.panel {
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