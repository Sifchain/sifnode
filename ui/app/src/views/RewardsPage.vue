<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import Loader from "@/components/shared/Loader.vue";

export default defineComponent({
  components: {
    Layout,
    Loader
  },
  data() {
    return {
      liqvalrewards: null,
    };
  },
  async mounted() {
    const data = await fetch(
      "https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/liqvalrewards"
    );
    const json = await data.json();
    this.liqvalrewards = json.body;
  },
});
</script>

<template>
  <Layout :header="false" title="Staking & Rewards" backLink="/peg">
    <div class="liquidity-container">
      <Loader black v-if="!liqvalrewards"/>
      <div v-else>
        <p class="mb-8">
          Earn additional ROWAN by staking or delegating! 
          The amount of rewards you can earn are: 
          <span v-if="liqvalrewards.liqValRewards === ''">TBD</span>
          <span v-else>
            {{liqvalrewards}} 
          </span>
          + Block rewards (variable)
        </p>
        <p>
          Learn more about staking and delegating <a href="https://docs.sifchain.finance/roles/validators" target="_blank">here</a>!
        </p>
      </div>
    </div>
  </Layout>
</template>

<style lang="scss">
.liquidity-container { 
  font-size: 18px;
  text-align: left;
  color: $c_gray_800;
  background: white; 
  padding: 15px !important
}
</style>