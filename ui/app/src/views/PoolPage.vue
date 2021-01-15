<script lang="ts">
import { defineComponent, ref } from "vue";
import { computed, toRefs } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { LiquidityProvider, Pool } from "ui-core";
import Layout from "@/components/layout/Layout.vue";
import PoolList from "@/components/poolList/PoolList.vue";
import PoolListItem from "@/components/poolList/PoolListItem.vue";
import SinglePool from "@/components/poolList/SinglePool.vue";
import SifButton from "@/components/shared/SifButton.vue";
type AccountPool = { lp: LiquidityProvider; pool: Pool };
export default defineComponent({
  components: {
    Layout,
    SifButton,
    PoolList,
    PoolListItem,
    SinglePool,
  },

  setup() {
    const { actions, poolFinder, store } = useCore();

    const selectedPool = ref<AccountPool | null>(null);
    const refsStore = toRefs(store);
    const accountPools = computed(() => refsStore.accountpools.value);

    return {
      accountPools,
      selectedPool,
    };
  },
});
</script>

<template>
  <SinglePool
    v-if="selectedPool"
    @back="selectedPool = null"
    :accountPool="selectedPool"
  />
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
        <PriceCalculation class="mb-8">
        <div class="info">
          <h3 class="mb-2">Liquidity provider rewards</h3>
          <p class="text--small mb-2">Liquidity providers earn a percentage fee on all trades proportional to their share of the pool. Fees are added to the pool, accrue in real time and can be claimed by withdrawing your liquidity. To learn more, reference of documentation <a href="https://docs.sifchain.finance/core-concepts/liquidity-pool">here</a></p>
        </div>
        </PriceCalculation>
        <PoolList class="mb-2">
            <PoolListItem
                    v-for="(accountPool, index) in accountPools"
                    :key="index"
                    :accountPool="accountPool"
                    @click="selectedPool = accountPool"
            />
      </PoolList>
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
