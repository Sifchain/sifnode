<template>
  <div class="pool-list-item">
    <div class="image">
      <img src="https://via.placeholder.com/22/0000FF" width="22" height="22">
      <img src="https://via.placeholder.com/22/00ff00" width="22" height="22">
    </div>
    <div class="symbol">
      <span>BTC</span>
      /
      <span>LSK</span>
    </div>
    <div class="button">
      Manage
    </div>
  </div>
</template>

<script> 
export default {

}
</script>

<style scoped lang="scss">
.pool-list-item {
  padding: 14px 16px;
  display: flex;

  &:not(:last-of-type) {
    border-bottom: $divider;
  }

  &:hover {
    cursor: pointer;
    background: $c_gray_50;
  }

  .image {
    height: 22px;

    img {
      border-radius: 16px;

      &:nth-child(2) {
        position: relative; 
        left: -6px;
      }
    }
  }

  .symbol {
    font-size: $fs_md;
    color: $c_text;
  }

  .button {
    font-size: $fs_sm;
    font-weight: normal;
    flex-grow: 1;
    text-align: right;
  }
}
</style>