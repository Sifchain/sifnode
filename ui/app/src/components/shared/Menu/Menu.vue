<template>
  <div class="menu">
    <div>
      <router-link class="logo" to="/">S</router-link>
    </div>
    <div class="item">
      <a
        class="link"
        target="_blank"
        href="https://docs.sifchain.finance/resources/sifchain-dex-ui"
        >Discover</a
      >
    </div>
    <div class="item">
      <a
        class="link"
        target="_blank"
        href="https://docs.sifchain.finance/resources/faq#how-can-i-acquire-rowan-and-or-erowan"
        >Get Rowan</a
      >
    </div>
    <div class="item">
      <router-link class="link" to="/stats">Pool Stats</router-link>
    </div>
    <div class="item">
      <router-link class="link" to="/stake-delegate"
        >Stake & Delegate</router-link
      >
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
    color: #c1a04f;
  }

  .link {
    @include resetLink;
    font-size: $fs_lg;
    font-weight: normal;
    opacity: 1;
    transition: opacity $trans_fast;
    color: #c1a04f;
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
