<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";


type AccountPool = { lp: LiquidityProvider; pool: Pool };

export default defineComponent({
  components: {
    Layout
  },

  setup() {
  },
});
</script>

<template>
  <Layout :header="true" title="Rewards"  backLink="/peg">
    <div class="liquidity-container">
      <div class="info">
        <p class="text--small mb-2">
          Rewards from Liquidity Mining Program: 49 Rowan YTD
        </p>
      </div>
    </div>
  </Layout>
</template>

<style scoped lang="scss">
.info {
    text-align: left;
}
.ul {text-decoration: underline}
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
