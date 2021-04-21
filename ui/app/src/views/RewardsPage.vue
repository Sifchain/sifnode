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
import Tooltip from "@/components/shared/Tooltip.vue";
import Icon from "@/components/shared/Icon.vue";
import DetailItem from "@/components/rewardsPage/DetailItem.vue";

const REWARD_INFO = {
  lm: {
    label: "Liquidity Minining",
    description:
      "Earn additional rewards by providing liquidity to any of Sifchain's pools.",
  },
  vs: {
    label: "Validator Subsidy",
    description:
      "Earn additional rewards by staking a node or delegating to a staked node.",
  },
};

const REWARD_DETAIL = {
  multiplier: {
    title: "Current Multiplier",
    tooltip:
      "Your multiplier as determined based on the amount of time you have held your position. At 121 days, this would be 4x.",
    icon: null,
    type: "number",
  },
  start: {
    title: "Start Date",
    tooltip: "start date",
    icon: null,
    type: "date",
  },
  reserved: {
    title: "Reserved Amount",
    tooltip:
      "The total amount of reward that has been reserved for you up to this point (based on the daily incremented reserved reward).",
    icon: "Rowan",
    type: "number",
  },
};

async function getRewardsData(address: ComputedRef<any>) {
  if (!address.value) return;
  const data = await fetch(
    `https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/rewards/${address.value}`,
  );
  if (data.status !== 200)
    return [
      { type: "lm", multiplier: 0, start: "", amount: null },
      { type: "vs", multiplier: 0, start: "", amount: null },
    ];

  return await data.json();
}

// NOTE - This will be removed and replaced with Amount API
function format(amount: number) {
  if (amount < 1) {
    return amount.toFixed(6);
  } else if (amount < 1000) {
    return amount.toFixed(4);
  } else if (amount < 100000) {
    return amount.toFixed(2);
  } else {
    return amount.toFixed(0);
  }
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
    Tooltip,
    Icon,
    DetailItem,
  },
  data() {
    return {
      showDetail: false,
    };
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
      format,
      REWARD_DETAIL,
    };
  },
});
</script>

<template>
  <Layout :header="true" title="Rewards">
    <Copy>
      Earn rewards by participating in any of our rewards-earning programs.
      Please see additional information of our
      <a
        target="_blank"
        href="https://docs.sifchain.finance/resources/rewards-programs"
        >current rewards program</a
      >
      and how to become eligible.
    </Copy>
    <div class="rewards-container">
      <div v-if="!rewards || rewards.length === 0" class="loader-container">
        <div class="loader" />
      </div>
      <Box v-else v-for="reward in rewards" v-bind:key="reward.type">
        <div class="reward-container">
          <div class="df fdr jcsb">
            <SubHeading>{{ REWARD_INFO[reward.type].label }}</SubHeading>
            <img
              class="cp"
              v-if="reward.detail"
              :class="reward.showDetail ? 'rotate-90' : ''"
              @click="reward.showDetail = !reward.showDetail"
              src="../assets/r-arrow.svg"
            />
          </div>
          <Copy>
            {{ REWARD_INFO[reward.type].description }}
          </Copy>
          <div class="details-container">
            <div class="amount-container">
              <div class="df fdr detail-item-amount" style="">
                <AssetItem symbol="Rowan" :label="false" />
                <span>{{ format(+reward.amount) }}</span>
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
          <div v-show="reward.showDetail" class="details-expanded-container">
            <DetailItem
              v-for="item in reward.detail"
              :key="Object.getOwnPropertyNames(item)[0]"
              :pkey="Object.getOwnPropertyNames(item)[0]"
              :item="item"
              :copy-map="REWARD_DETAIL"
            />
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
.detail-item-amount {
  width: 120px;
}
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
