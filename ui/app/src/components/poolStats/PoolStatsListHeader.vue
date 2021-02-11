<template>
  <div class="header">
    <div class="col-sm">
      <span>Token</span>
    </div>
    <div class="col-md">
      <span>Value of Pool</span>
      <span>(USD)</span>
    </div>
    <div class="col-md">
      <span>Value of Pool</span>
      <span>(ROWAN)</span>
    </div>
    <div class="col-md">
      <span>Value of External Token</span>
      <span>(USD)</span>
    </div>
    <div class="col-md">
      <span>Value of External Token</span>
      <span>(ROWAN)</span>
    </div>
    <div class="col-md">
      <span>Liquidity Mining Bonus APY</span>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .header {
    padding: 12px 12px;
    background: $c_gray_900;
    display: flex;
    justify-content: space-evenly;
    align-items: flex-start;

    &:hover {
      cursor: pointer;
      background: $c_gray_800;
    }
  }

  .col-sm {
    width: 10%;
    color: $c_gray_400;
    justify-items: start;
  }

  .col-md {
    width: 18%;
    color: $c_gray_400;
  }
</style>
