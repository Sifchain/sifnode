<template>
  <div class="pool-units">
    <div class="pool-units-value mb-10">
      {{ poolUnits }}
    </div>
    <div class="pool-asset">
      <AssetItemPool 
        :tokenASymbol="fromSymbol"
        :tokenBSymbol="toSymbol"
        inline
      />
      <span>&nbsp;Pool Tokens</span>
    </div>
    <!-- <div class="pool-units-label">
      {{ fromSymbol }}/{{ toSymbol }} Pool Tokens<br />
    </div> -->
  </div>
</template>

<script>
import AssetItemPool from "@/components/shared/AssetItemPool.vue";

export default {
  components: {
    AssetItemPool,
  },

  props: {
    poolUnits: String,
    fromSymbol: String,
    toSymbol: String,
  }
}
</script>

<style lang="scss" scoped>
.pool-units {
  margin-bottom: 32px;

  &-value {
    font-size: 24px;
    font-style: normal !important;
    text-align: left;
  }
  
  &-label {
    text-align: left;
    font-weight: 400;
  }
}
.pool-asset {
  display: flex;
  align-items: center;
}
</style>