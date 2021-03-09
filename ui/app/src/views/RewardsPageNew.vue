<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import SifButton from "@/components/shared/SifButton.vue";

export default defineComponent({
  components: {
    Layout,
    SifButton,
  },

  setup() {},
});
</script>

<template>
  <Layout :header="true" title="Rewards" backLink="/peg">
    <div class="liquidity-container">
      <div class="info">
        <p class="mb-5">
          Rewards from Liquidity Mining Program:
        </p>
        <p class="mb-2">
          <span class="amount">49 Rowan YTD</span>
        </p>
        <SifButton disabled primary class="mb-8">Claim</SifButton>
        <p class="text--small">
          At present, you are unable to claim Liquidity Mining rewards.
        </p>
      </div>
    </div>
  </Layout>
</template>

<style scoped lang="scss">
.info {
  text-align: left;
}
.ul {
  text-decoration: underline;
}
.amount {font-size: $fs_lg}
.liquidity-container {
  text-align: left;
  color: $c_gray_700;
  border-top: 1px solid $c_gray_400;
  min-height: 145px;
  background: white;
  padding: 15px;
  border-radius: 0 0 6px 6px;
  padding-bottom: 60px
}
</style>
