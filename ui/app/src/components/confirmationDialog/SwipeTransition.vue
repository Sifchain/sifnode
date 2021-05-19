<template>
  <transition name="swipe">
    <slot></slot>
  </transition>
</template>
<style lang="scss">
.swipe-enter-active,
.swipe-leave-active {
  transition: transform 0.5s ease-out;
}
.swipe-enter-from {
  transform: translateX(100%);
}
.swipe-leave-to {
  transform: translateX(-100%);
}
</style>
