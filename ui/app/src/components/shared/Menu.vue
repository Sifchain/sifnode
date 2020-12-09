<template>
  <div class="menu">
    <div>
      <router-link class="logo" to="/">S</router-link>
    </div>
    <div class="item">
      <router-link class="link" to="/">Discover</router-link>
    </div>
    <div class="item">
      <router-link class="link" to="/">New wallet</router-link>
    </div>
    <div class="item">
      <router-link class="link" to="/">Convert</router-link>
    </div>
    <div class="item">
      <router-link class="link" to="/">Deposit</router-link>
    </div>
  </div>
</template>

<script>
export default {};
</script>

<style lang="scss" scoped>
.menu {
  display: flex;
  align-items: baseline;
  .item {
    margin-right: 35px;
  }

  .logo {
    @include resetLink;
    font-family: $f_special;
    font-size: 40px;
    font-weight: normal;
    margin-right: 46px;
  }

  .link {
    @include resetLink;
    font-size: $fs_lg;
    font-weight: normal;
    opacity: 1;
    transition: opacity $trans_fast;

    &:hover {
      opacity: 0.5;
    }

    &.active {
      font-weight: bold;
      pointer-events: none;
    }
  }
}
</style>