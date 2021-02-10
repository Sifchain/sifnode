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
    poolData: {
      type: Object,
    },
    inline: Boolean,
  },
});
</script>

<template>
  <div class="pool-list">
    <PoolStatsListHeader />
    <PoolStatsListItem
      v-if="poolData && poolData.pools && poolData.liqAPY && poolData.pools[0]"
      v-for="(pool, index) in poolData.pools"
      :key="index"
      :pool="pool"
      :liqAPY="poolData.liqAPY"
    />
    <div v-else class="loading">
      <div class="ring"></div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.pool-list {
  overflow-y: auto;
  background: $c_white;
  border-radius: $br_sm;
}

.loading {
  height: 240px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 20px;
}

.ring {
  position: absolute;
  top: 300px;
  width: 64px;
  height: 64px;
  border: 4px solid $c_gold;
  border-radius: 40px;
  border-top-color: transparent;
  animation: 1s rotate infinite linear;
}

@keyframes rotate {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
