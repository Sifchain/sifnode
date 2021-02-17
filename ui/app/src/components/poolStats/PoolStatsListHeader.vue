<template>
  <div class="header">
    <div class="col-sm">
      <span>Token</span>
    </div>
    <div class="col-sm">
      <span>Price of Token</span>
      <span>(USD)</span>
    </div>
    <div class="col-sm">
      <span>Pool Depth</span>
      <span>(USD)</span>
    </div>
    <div class="col-sm">
      <span>Trade Volume</span>
      <span>(24h)</span>
    </div>
    <div class="col-sm">
      <span>Pool APY</span>
    </div>
    <div class="col-lg">
      <span>Liquidity Mining Bonus</span>
      <span>APY</span>
    </div>
    <div class="col-sm">
      <span>Total APY</span>
    </div>
  </div>
</template>

<style scoped lang="scss">
.header {
  padding: 2em 1em;
  min-height: 64px;
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
  min-width: 102px;
  color: $c_gray_400;
  display: flex;
  flex-direction: column;
}

.col-md {
  min-width: 110px;
  color: $c_gray_400;
  display: flex;
  flex-direction: column;
}
.col-lg {
  min-width: 168px;
  color: $c_gray_400;
  display: flex;
  flex-direction: column;
}
</style>
