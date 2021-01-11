<script lang="ts">
import { defineComponent, ref } from "vue";
import { computed, reactive  } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { Pool } from "ui-core";
import Layout from "@/components/layout/Layout.vue";
import PoolList from "@/components/poolList/PoolList.vue";
import PoolListItem from "@/components/poolList/PoolListItem.vue";
import SinglePool from "@/components/poolList/SinglePool.vue";
import SifButton from "@/components/shared/SifButton.vue";
import PriceCalculation from "@/components/shared/PriceCalculation.vue";

export default defineComponent({
  components: { Layout, SifButton, PriceCalculation, PoolList, PoolListItem, SinglePool },
  
  setup() {
    const { actions, poolFinder, store } = useCore();
    const pools = ref<any>(null);
    const selectedPool = ref<any>(null);

    async function getPools() {
      await actions.clp.getLiquidityProviderPools().then((res)=>{
        pools.value = res
      })
    }

    getPools();

    return {
      pools,
      selectedPool,
    }
  }
});
</script>

<template>
  <SinglePool v-if="selectedPool" @back="selectedPool = null" :pool="selectedPool" />
  <Layout v-else>
    <div>
      <div class="heading mb-8">
        <h3>Your Liquidity</h3>
        <router-link to="/pool/create-pool" class="pr-4"
          ><SifButton primaryOutline nocase>Create Pair</SifButton></router-link
        >&nbsp;
        <router-link to="/pool/add-liquidity"
          ><SifButton primary nocase>Add Liquidity</SifButton></router-link
        >
      </div>
      <div class="mb-8">
        <SifButton primaryOutline nocase block>Account analytics and accrued fees</SifButton>
      </div>
      <PriceCalculation class="mb-8">
        <div class="info">
          <h3 class="mb-2">Liquidity provider rewards</h3>
          <p class="text--small mb-2">Liquidity providers earn a 0.3% fee on all trades proportional to their share of the pool. Fees are added to the pool, accrue in real time and can be claimed by withdrawing your liquidity.</p>
          <p class="text--small mb-2"><a href="#">Read more about providing liquidity</a></p>
        </div>
      </PriceCalculation>
      <PoolList class="mb-2">
        <PoolListItem v-for="(pool, index) in pools" :key="index" :pool="pool" @click="selectedPool = pool"/>
      </PoolList>
      <div class="footer">
        Don’t see a pool you joined? <a href="#">Import it</a>
      </div>
    </div>
  </Layout>
</template>

<style scoped lang="scss">
.heading {
  display: flex;
  align-items: center;

  h3 {
    @include title16;
  }
}

.info {
  text-align: left;
  padding: 8px;
  font-weight: 400;
}

.footer {
  font-weight: 400;
}
</style>