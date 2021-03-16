<template>
  <i class="caret-arrow caret-down" />
</template>
<style lang="scss" scoped>
/* FROM https://www.w3schools.com/howto/howto_css_arrows.asp make props work for others */
.caret-arrow {
  border: solid $c_gray_700;
  border-width: 0 2px 2px 0;
  display: inline-block;
  padding: 2px;
}

.caret-right {
  transform: rotate(-45deg);
}

.caret-left {
  transform: rotate(135deg);
}

.caret-up {
  transform: rotate(-135deg);
}

.caret-down {
  position: relative;
  top: -4px;
  transform: rotate(45deg);
}
</style>
