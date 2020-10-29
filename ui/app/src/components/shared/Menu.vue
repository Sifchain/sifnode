<template>
  <div class="menu">
    <div class="menu__item">
      <a class="menu__link menu__link--active" href="#">Discover</a>
    </div>
    <div class="menu__item">
      <a class="menu__link" href="#">New Wallet</a>
    </div>
    <div class="menu__item">
      <a class="menu__link" href="#">Convert</a>
    </div>
    <div class="menu__item">
      <a class="menu__link" href="#">Deposit</a>
    </div>
  </div>
</template>

<script>
export default {

}
</script>

<style lang="scss">
.menu {
  display: flex;

  &__item {
    margin-right: 35px;
  }

  &__link {
    @include resetLink;
    font-size: $fs_lg;
    font-weight: normal;
    opacity: 1;
    transition: opacity $trans_fast;

    &:hover {
      opacity: 0.5;
    }

    &--active {
      font-weight: bold;
      pointer-events: none;
    }
  }
}
</style>