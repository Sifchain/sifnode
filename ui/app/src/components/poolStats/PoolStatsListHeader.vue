<template>
  <div class="header">
    <div class="col-sm">
      <span>Token</span>
    </div>
    <div class="col-md">
      <span>Pool Depth in USD</span>
    </div>
    <div class="col-md">
      <span>Pool Depth in Rowan</span>
    </div>
    <div class="col-lg">
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
  min-width: 100px;
  width: 19%;
  color: $c_gray_400;
  justify-items: start;
}

.col-md {
  min-width: 160px;
  width: 27%;
  color: $c_gray_400;
}
.col-lg {
  min-width: 200px;
  width: 27%;
  color: $c_gray_400;
}
</style>
