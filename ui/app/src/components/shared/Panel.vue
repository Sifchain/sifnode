<template>
  <div :class="{ panel: true, dark }">
    <div class="header" v-if="!!$slots.header">
      <slot name="header"></slot>
    </div>

    <div class="body">
      <slot></slot>
    </div>

    <div class="footer" v-if="!!$slots.header">
      <slot name="footer"></slot>
    </div>
  </div>
  <!-- main app card -->
</template>
<script>
export default {
  props: {
    dark: Boolean,
  },
};
</script>
<style lang="scss" scoped>
.panel::-webkit-scrollbar {
  display: none;
}
.panel {
  width: 410px;
  border-radius: $br_lg;
  overflow: scroll;
  background: $c_white;
  position: relative;

  .header {
    padding: 20px 20px 10px 20px;
  }

  &.dark {
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
}
</style>