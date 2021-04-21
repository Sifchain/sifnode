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
        >DISCOVER</a
      >
    </div>
    <div class="item">
      <a
        class="link"
        target="_blank"
        href="https://docs.sifchain.finance/resources/faq#how-can-i-acquire-rowan-and-or-erowan"
        >GET ROWAN</a
      >
    </div>
    <div class="item">
      <router-link class="link" to="/stats">POOL STATS</router-link>
    </div>
    <div class="item">
      <router-link class="link" to="/stake-delegate"
        >STAKE & DELEGATE</router-link
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
    font-size: 14px;
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
