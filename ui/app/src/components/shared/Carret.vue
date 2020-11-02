<template>
  <i class="carret-arrow carret-down" />
</template>
<style lang="scss" scoped>
/* FROM https://www.w3schools.com/howto/howto_css_arrows.asp make props work for others */
.carret-arrow {
  border: solid $c_gray_700;
  border-width: 0 2px 2px 0;
  display: inline-block;
  padding: 2px;
}

.carret-right {
  transform: rotate(-45deg);
}

.carret-left {
  transform: rotate(135deg);
}

.carret-up {
  transform: rotate(-135deg);
}

.carret-down {
  position: relative;
  top: -4px;
  transform: rotate(45deg);
}
</style>