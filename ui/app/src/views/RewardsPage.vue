<script lang="ts">
import { computed, defineComponent, watchEffect } from "vue";
import { ref } from "@vue/reactivity";
import { Amount, IAmount } from "ui-core";
import { useCore } from "@/hooks/useCore";
import Layout from "@/components/layout/Layout.vue";
import SifButton from "@/components/shared/SifButton.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import Box from "@/components/shared/Box.vue";
import { Copy, SubHeading } from "@/components/shared/Text";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import Tooltip from "@/components/shared/Tooltip.vue";
import Icon from "@/components/shared/Icon.vue";
import { format } from "ui-core/src/utils/format";

const REWARD_INFO = {
  lm: {
    label: "Liquidity Mining",
    description:
      "Earn additional rewards by providing liquidity to any of Sifchain's pools.",
  },
};

// Need verification this is the correct schema
type RewardsResult = {
  type: string;
  multiplier: IAmount;
  start: string; // date string
  amount: IAmount;
};

function toAmount(thing: any) {
  try {
    return Amount(thing.toString());
  } catch (err) {
    return Amount("0");
  }
}

// TODO: this cannot be tested properly right now we need to provide this as an injectable service
// through a usecase to adequately test the data We might want to add this service to our docker
// backing stack or use fixtures within e2e tests we will likely want to add wrap this around some
// kind of swr pattern to have a better UX
async function getRewardsData(address: string): Promise<RewardsResult[]> {
  if (!address) return [];

  const data = await fetch(
    `https://vtdbgplqd6.execute-api.us-west-2.amazonaws.com/default/rewards/${address}`,
  );

  if (data.status !== 200)
    return [
      { type: "lm", multiplier: Amount("0"), start: "", amount: Amount("0") },
    ];

  const rewardsRaw = (await data.json()) as RewardsResult[];

  // Map to Amount API
  return rewardsRaw.map(
    ({ type, multiplier, start, amount }: RewardsResult) => {
      return {
        type,
        multiplier: toAmount(multiplier),
        start,
        amount: toAmount(amount),
      };
    },
  );
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
  },
  setup() {
    const { store } = useCore();
    const address = computed(() => store.wallet.sif.address);
    let rewards = ref<RewardsResult[]>([]);

    watchEffect(async () => {
      if (address.value) {
        const data = await getRewardsData(address.value);
        console.log({ data });
        rewards.value = data;
      }
    });

    // TODO: save this somewhere as a reuse option
    const dynamicMantissa = {
      1: 6,
      1000: 4,
      100000: 2,
      infinity: 0,
    };

    return {
      rewards,
      REWARD_INFO,
      format,
      dynamicMantissa,
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
          <SubHeading>{{ REWARD_INFO[reward.type].label }}</SubHeading>
          <Copy>
            {{ REWARD_INFO[reward.type].description }}
          </Copy>
          <div class="details-container">
            <div class="amount-container w50 jcsb">
              <div class="df fdr">
                <AssetItem symbol="Rowan" :label="false" />
                <span data-handle="reward-amount">{{
                  format(reward.amount, { mantissa: dynamicMantissa })
                }}</span>
              </div>
              <span>ROWAN</span>
              <Tooltip>
                <template #message>
                  <div class="tooltip">
                    Current multiplier:
                    {{
                      format(reward.multiplier, { mantissa: dynamicMantissa })
                    }}x
                  </div>
                </template>
                <Icon icon="info-box-black" data-handle="multiply-tooltip" />
              </Tooltip>
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
