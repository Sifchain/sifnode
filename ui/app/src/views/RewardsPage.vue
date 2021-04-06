<script lang="ts">
import { computed, defineComponent, watch, onMounted } from "vue";
import { ref, ComputedRef } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import Layout from "@/components/layout/Layout.vue";
import SifButton from "@/components/shared/SifButton.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import Box from "@/components/shared/Box.vue";
import { Copy, SubHeading } from "@/components/shared/Text";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";

const REWARD_INFO = {
  lm: {
    label: "Liquidity Minining",
    description:
      "Earn additional rewards by providing liquidity to any of Sifchain's pools.",
  },
};

async function getRewardsData(address: ComputedRef<any>) {
  if (!address.value) return;
  const data = await fetch(
    `https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/rewards/${address.value}`,
  );
  if (data.status !== 200)
    return [{ type: "lm", multiplier: 0, start: "", amount: null }];
  return await data.json();
}
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
    let rewards = ref<Array<Object>>([]);

    watch(address, async () => {
      rewards.value = await getRewardsData(address);
    });

    onMounted(async () => {
      rewards.value = await getRewardsData(address);
    });

    return {
      rewards,
      REWARD_INFO,
    };
  },
});
</script>

<template>
  <Layout :header="true" title="Rewards">
    <Copy>
      Earn rewards by participating in of our rewards-earning programs. Please
      see additional information of our
      <a
        target="_blank"
        href="https://docs.sifchain.finance/resources/rewards-programs"
        >current rewards program</a
      >
      and how to become eligible.
    </Copy>
    <div class="rewards-container">
      <div v-if="rewards.length === 0" class="loader-container">
        <div class="loader" />
      </div>
      <Box v-else v-for="reward in rewards" v-bind:key="reward.type">
        <div class="reward-container">
          <SubHeading>{{ REWARD_INFO[reward.type].label }}</SubHeading>
          <Copy>
            Earn additional rewards by staking a node or delegating to a staked
            node.
          </Copy>
          <div class="details-container">
            <div class="amount-container w50 jcsb">
              <div class="df fdr">
                <AssetItem symbol="Rowan" :label="false" />
                <span>{{ reward.amount ? reward.amount?.toFixed() : 0 }}</span>
              </div>
              <span>ROWAN</span>
            </div>
            <a
              class="more-info-button"
              target="_blank"
              href="https://docs.sifchain.finance/resources/rewards-programs#liquidity-mining-and-validator-subsidy-rewards-on-sifchain"
              >More Info</a
            >
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
.more-info-button {
  // TODO - This Button !
  background: #f3f3f3;
  color: #343434;
  font-size: 12px;
  border-radius: 6px;
  width: 96px;
  height: 30px;
  font-weight: 100;
  display: flex;
  justify-content: center;
  align-items: center;
}

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
    > :nth-child(1),
    > :nth-child(2) {
      margin-bottom: $margin_small;
    }
    .amount-container {
      display: flex;
      align-items: center;
      flex-direction: row;
    }
    .details-container {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
    }
  }

  /* TODO - TEMP - Need to componentize our loaders */
  .loader-container {
    margin-top: $margin-large;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .loader {
    background: url("../../public/images/siflogo.png");
    background-size: cover;
    width: 64px;
    height: 64px;
    box-shadow: 0 0 0 0 rgba(0, 0, 0, 1);
    transform: scale(1);
    animation: pulse 1s infinite;
  }
  @keyframes pulse {
    0% {
      transform: scale(0.85);
    }
    70% {
      transform: scale(1);
    }
    100% {
      transform: scale(0.85);
    }
  }
}
</style>
