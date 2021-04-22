<template>
  <div :class="{ panel: true, dark }">
    <div
      class="header"
      v-if="!!$slots.header"
      :data-handle="[dataHandle, 'panel-header'].filter(Boolean).join('-')"
    >
      <slot name="header"></slot>
    </div>

    <div
      class="body"
      :data-handle="[dataHandle, 'panel-body'].filter(Boolean).join('-')"
    >
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
    dataHandle: String,
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
