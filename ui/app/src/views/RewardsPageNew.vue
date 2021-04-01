<script lang="ts">
import { computed, defineComponent, watch } from "vue";
import Layout from "@/components/layout/Layout.vue";
import SifButton from "@/components/shared/SifButton.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import { useCore } from "@/hooks/useCore";
import { ref } from "@vue/reactivity";

export default defineComponent({
  components: {
    Layout,
    SifButton,
    AssetItem,
    ActionsPanel
  },
  setup() {
    const { store } = useCore();
    const address = computed(() => store.wallet.sif.address);
    let rewards = ref<Array<Object>>([])

    watch(
      address, 
      async () => {
        const data = await fetch(`https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/rewards/${address.value}`);
        rewards.value = await data.json();
      }
    )
    return {
      rewards
    }
  },
});
</script>

<template>
  <Layout :header="true" title="Rewards" backLink="/peg">
    <div class="info mb-8" style="line-height: 1.4">
      <h3 class="mb-2">Your Rewards</h3>
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
    <div class="list-container" 
      v-for="reward in rewards" 
      v-bind:key="reward.type"
    >
      <div class="item">
        <div class="title">
          {{reward.type}}
        </div>

        <div class="detail">
          <!-- future: slot -->
          <div class="amount">
            <AssetItem symbol="Rowan" :label="false" />
            <span class="mr-6">{{reward.amount}}</span>
          </div>
          <SifButton primary>Claim</SifButton>
        </div>
      </div>
    </div>
    <ActionsPanel
      connectType="connectToSif"
    />
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
