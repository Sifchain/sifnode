<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import Loader from "@/components/shared/Loader.vue";

export default defineComponent({
  components: {
    Layout,
    Loader,
  },
  data() {
    return {
      data: null,
    };
  },
  async mounted() {
    const data = await fetch(
      "https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/liqvalrewards",
    );
    const json = await data.json();
    this.data = json.body;
  },
});
</script>

<template>
  <Layout :header="false" title="Staking & Rewards" backLink="/balances">
    <div class="liquidity-container">
      <Loader black v-if="!data" />
      <div v-else>
        <p class="mb-8">
          Earn additional Rowan by staking or delegating. The amount of rewards
          you can earn is the summation of:
        </p>
        <p class="">
          1.
          <a
            class="ul"
            href="https://medium.com/sifchain-finance/uses-for-rowan-the-polyvalent-token-for-omni-chain-decentralized-exchange-dex-3207e7f70f02?source=collection_home---4------10-----------------------"
            target="_blank"
            >Validator Subsidy</a
          >:
          <span v-if="data.liqValRewards === ''">TBD</span>
          <span v-else>
            {{ parseFloat(data.liqValRewards).toFixed(2) }} % APY
          </span>
        </p>
        <p class="mb-8">
          2.
          <a
            class="ul"
            href="https://docs.sifchain.finance/roles/validators#block-rewards"
            target="_blank"
            >Block Rewards</a
          >: (variable)
        </p>
        <p class="mb-8">
          Learn more about staking and delegating
          <a
            class="ul"
            href="https://docs.sifchain.finance/roles/validators"
            target="_blank"
            >here</a
          >!
        </p>
        <p class="mb-9">
          Delegation instructions
          <a
            class="ul"
            href="https://docs.sifchain.finance/roles/delegators#how-to-delegate"
            target="_blank"
            >here</a
          >
        </p>
      </div>
    </div>
  </Layout>
</template>

<style scoped lang="scss">
.ul {
  text-decoration: underline;
}
.liquidity-container {
  text-align: left;
  color: $c_gray_700;
  border-top: 1px solid $c_gray_400;
  min-height: 145px;
  background: white;
  padding: 15px;
  border-radius: 0 0 6px 6px;
}
</style>
