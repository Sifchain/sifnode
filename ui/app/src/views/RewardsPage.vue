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
      data: null,
    };
  },
  async mounted() {
    const data = await fetch(
      "https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/liqvalrewards"
    );
    const json = await data.json();
    this.data = json.body;
  },
});
</script>

<template>
  <Layout :header="false" title="Staking & Rewards" backLink="/peg">
    <div class="liquidity-container">
      <Loader black v-if="!data"/>
      <div v-else>
        <p class="mb-8">
          Earn additional ROWAN by staking or delegating! 
          The amount of rewards you can earn are: 
          <span v-if="data.liqValRewards === ''">TBD</span>
          <span v-else>
            {{data.liqValRewards}} 
          </span>
          + Block rewards (variable)
        </p>
        <p class="mb-9">
          Learn more about staking and delegating <a href="https://docs.sifchain.finance/roles/validators" target="_blank">here</a>!
        </p>
      </div>
    </div>
  </Layout>
</template>

<style scoped lang="scss">
.liquidity-container { 
  text-align: left;
  color: $c_gray_700;
  border-top: 1px solid $c_gray_400;
  min-height: 145px;
  background: white; 
  padding: 15px;
  border-radius: 0 0 6px 6px
}
</style>