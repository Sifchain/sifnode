<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import SifButton from "@/components/shared/SifButton.vue";
import AssetItem from "@/components/shared/AssetItem.vue";

export default defineComponent({
  components: {
    Layout,
    SifButton,
    AssetItem,
  },
  setup() {},
});
</script>

<template>
  <Layout :header="true" title="Rewards" backLink="/peg">
    <div class="info mb-8" style="line-height: 1.4">
      <p class="text--small mb-2">
        Earn rewards by participating in of our rewards-earning programs. Please
        see additional information of our current rewards programs and how to
        become eligible for them
        <a
          target="_blank"
          href="https://docs.sifchain.finance/resources/rewards-programs"
          >here</a
        >.
      </p>
    </div>
    <div class="list-container">
      <div class="item">
        <div class="title">
          Liquidity Mining
        </div>

        <div class="detail">
          <!-- future: slot -->
          <div class="amount">
            <AssetItem symbol="Rowan" :label="false" />
            <span class="mr-6">49.212</span>
          </div>
          <SifButton primary>Claim</SifButton>
        </div>
      </div>
    </div>
  </Layout>
</template>

<style scoped lang="scss">
/* this pattern is generic. for component lib */
.info {
  text-align: left;
  font-weight: 400;
}
.list-container {
  text-align: left;
  color: $c_gray_700;
  border-top: 1px solid $c_gray_400;
  min-height: 145px;
  background: white;
  border-radius: 0 0 6px 6px;
  padding-bottom: 60px;
  .item {
    padding: 14px 16px;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    border-bottom: $divider;
    align-items: center;
    &:hover {
      cursor: pointer;
      background: $c_gray_50;
    }
    .title {
      color: $c_text;
    }
    .detail {
      display: flex;
      align-items: center;
      .amount {
        display: flex;
      }
    }
  }
}
</style>
