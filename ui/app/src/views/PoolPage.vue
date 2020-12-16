<script lang="ts">
import { defineComponent, ref } from "vue";
import { computed, reactive } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import Layout from "@/components/layout/Layout.vue";
import PoolList from "@/components/poolList/PoolList.vue";
import PoolListItem from "@/components/poolList/PoolListItem.vue";
import SifButton from "@/components/shared/SifButton.vue";
import PriceCalculation from "@/components/shared/PriceCalculation.vue";

export default defineComponent({
  components: { Layout, SifButton, PriceCalculation, PoolList, PoolListItem },

  mounted() {
    this.getPools()
  },
  
  setup() {
    const { actions, poolFinder, store } = useCore();
    const state = reactive({
      pools: ref(),
      selectedPool: null,
    })

    async function getPools() {
      await actions.clp.getLiquidityProviderPools().then((res) => {
        state.pools = res
      })
    }

    getPools();

    return {
      getPools: getPools,
      state,
      poolSelected: (index:number) => {
        state.selectedPool = state.pools[index]
        console.log(state.selectedPool)
      }
    }
  }
});
</script>

<template>
  <Layout>
    
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
        <PoolListItem v-for="(pool, index) in state.pools" :key="index" :pool="pool" @poolSelected="poolSelected(index)"/>
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