<template>
  <header class="header">
    <div class="header__left">
      <div class="header__logo">S</div>
      <Menu />
    </div>
    <div class="header__right">
      <SifButton>Rowan</SifButton>
    </div>
  </header>
</template>

<script>
import Menu from './Menu.vue';
import SifButton from '../Elements/SifButton.vue';

export default {
  components: {
    Menu,
    SifButton,
  }
}
</script>

<style lang="scss">
.header {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 80px;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  border-bottom: $divider;

  &__left {
    display: flex;

  }

  &__logo {
    font-family: $f_special;
    font-size: 40px;
    font-weight: normal;
    margin-right: 46px;
  }
}
</style>