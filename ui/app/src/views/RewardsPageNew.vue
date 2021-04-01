<script lang="ts">
import { computed, defineComponent, watch } from "vue";
import { ref } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import Layout from "@/components/layout/Layout.vue";
import SifButton from "@/components/shared/SifButton.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import Box from "@/components/shared/Box.vue";
import { Copy, SubHeading } from "@/components/shared/Text";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
export default defineComponent({
  components: {
    Layout,
    SifButton,
    AssetItem,
    ActionsPanel,
    Copy,
    SubHeading,
    Box,
  },
  setup() {
    const { store } = useCore();
    const address = computed(() => store.wallet.sif.address);
    let rewards = ref<Array<Object>>([
      { type: "lm", multiplier: 0, start: "", amount: 12709.098861115788 },
      { type: "lm", multiplier: 0, start: "", amount: 333.098861115788 },
    ]);
    //
    // watch(address, async () => {
    //   const data = await fetch(
    //     `https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/rewards/${address.value}`,
    //   );
    //   rewards.value = await data.json();
    // });
    return {
      rewards,
    };
  },
});
</script>

<template>
  <Layout :header="true" title="Rewards" backLink="/peg">
    <Copy>
      Earn rewards by participating in of our rewards-earning programs. Please
      see additional information of our current rewards programs and how to
      become eligible for them
      <a
        target="_blank"
        href="https://docs.sifchain.finance/resources/rewards-programs"
        >here</a
      >.
    </Copy>
    <div class="rewards-container">
      <Box v-for="reward in rewards" v-bind:key="reward.type">
        <div class="reward-container">
          <SubHeading>{{ reward.type }}</SubHeading>
          <div class="details-container">
            <div class="amount-container">
              <AssetItem symbol="Rowan" :label="false" />
              <span>{{ reward.amount }}</span>
            </div>
            <SifButton primary>Claim</SifButton>
          </div>
        </div>
      </Box>
    </div>
    <ActionsPanel connectType="connectToSif" />
  </Layout>
</template>

<style scoped lang="scss">
// TODO - Get variable margin/padding sizes in
// TODO - Discuss how we should manage positioning

.rewards-container {
  display: flex;
  flex-direction: column;
  > :first-child {
    margin-top: $margin_medium;
  }
  width: 100%;
  > :nth-child(1) {
    margin-bottom: $margin_medium;
  }
  .reward-container {
    flex-direction: column;
    .amount-container {
      display: flex;
      flex-direction: row;
    }
    .details-container {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
    }
  }
}

/*
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
*/
</style>
