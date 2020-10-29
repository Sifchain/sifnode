<template>
  <div class="menu">
    <div class="item">
      <a class="link active" href="#">Discover</a>
    </div>
    <div class="item">
      <a class="link" href="#">New Wallet</a>
    </div>
    <div class="item">
      <a class="link" href="#">Convert</a>
    </div>
    <div class="item">
      <a class="link" href="#">Deposit</a>
    </div>
  </div>
</template>

<script>
export default {

}
</script>

<style lang="scss" scoped>
.menu {
  display: flex;

  .item {
    margin-right: 35px;
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