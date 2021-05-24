<script>
import { defineComponent } from "vue";
import PoolStatsListItem from "@/components/poolStats/PoolStatsListItem.vue";
import PoolStatsListHeader from "@/components/poolStats/PoolStatsListHeader.vue";

export default defineComponent({
  components: {
    PoolStatsListItem,
    PoolStatsListHeader,
  },

  props: {
    liqAPY: {
      type: Number,
    },
    poolData: {
      type: Object,
    },
    inline: Boolean,
  },
});
</script>

<template>
  <PoolStatsListHeader
    v-if="poolData && poolData.pools && poolData.liqAPY && poolData.pools[0]"
    class="pool-list-header"
  />
  <div class="pool-list-container">
    <PoolStatsListItem
      v-if="poolData && poolData.pools && poolData.liqAPY && poolData.pools[0]"
      v-for="(pool, index) in poolData.pools"
      :key="index"
      :pool="pool"
      :liqAPY="liqAPY"
      class="pool-list"
    />
    <div v-else class="loading">
      <div class="logo"></div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.pool-list-header {
  border-top-left-radius: $br_sm;
  border-top-right-radius: $br_sm;
}

.pool-list-container {
  overflow-y: auto;
  overflow-x: hidden;
  margin-bottom: 32px;
  border-bottom-left-radius: $br_sm;
  border-bottom-right-radius: $br_sm;
}

.pool-list {
  background: $c_white;
}

.loading {
  margin-top: 180px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.logo {
  background: url("../../../public/images/siflogo.png");
  background-size: cover;
  width: 64px;
  height: 64px;
  box-shadow: 0 0 0 0 rgba(0, 0, 0, 1);
  transform: scale(1);
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(0.85);
  }
  70% {
    transform: scale(1);
  }
  100% {
    transform: scale(0.85);
  }
}
</style>
