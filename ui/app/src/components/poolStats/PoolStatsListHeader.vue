<template>
  <div class="header">
    <div class="col-sm">
      <span>Token</span>
    </div>
    <div class="col-md">
      <span>Price of Token</span>
      <span>(USD)</span>
    </div>
    <div class="col-md">
      <span>Pool Depth</span>
      <span>(USD)</span>
    </div>
    <div class="col-sm">
      <span>Volume</span>
      <span>(24h)</span>
    </div>
    <div class="col-lg">
      <span>Liquidity Mining Bonus</span>
      <span>APY</span>
    </div>
  </div>
</template>

<style scoped lang="scss">
  .header {
    padding: 12px 12px;
    background: $c_gray_900;
    display: flex;
    justify-content: space-evenly;
    align-items: center;

    &:hover {
      cursor: pointer;
      background: $c_gray_800;
    }
  }

  .col-sm {
    min-width: 100px;
    width: 15%;
    color: $c_gray_400;
    display: flex;
    flex-direction: column;
  }

  .col-md {
    min-width: 160px;
    width: 20%;
    color: $c_gray_400;
    display: flex;
    flex-direction: column;
  }
  .col-lg {
    min-width: 200px;
    width: 27%;
    color: $c_gray_400;
    display: flex;
    flex-direction: column;
  }
</style>
