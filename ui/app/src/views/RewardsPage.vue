<script lang="ts">
import { computed, defineComponent, watch, onMounted } from "vue";
import { ref, ComputedRef } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { getCryptoeconomicsUrl } from "@/components/shared/utils";
import Layout from "@/components/layout/Layout.vue";
import SifButton from "@/components/shared/SifButton.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import ConfirmationModal from "@/components/shared/ConfirmationModal.vue";
import Box from "@/components/shared/Box.vue";
import { Copy, SubHeading } from "@/components/shared/Text";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import Modal from "@/components/shared/Modal.vue";
import ModalView from "@/components/shared/ModalView.vue";
import PairTable from "@/components/shared/PairTable.vue";
import Tooltip from "@/components/shared/Tooltip.vue";
import Icon from "@/components/shared/Icon.vue";
import { ConfirmState } from "@/types";
const REWARD_INFO = {
  lm: {
    label: "Liquidity Mining",
    description:
      "Earn additional rewards by providing liquidity to any of Sifchain's pools.",
  },
  vs: {
    label: "Validator Subsidy",
    description:
      "Earn additional rewards by staking a node or delegating to a staked node.",
  },
};

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

async function getLMData(address: ComputedRef<any>, chainId: string) {
  if (!address.value) return;
  // const timestamp = Date.parse(new Date().toString());
  const ceUrl = getCryptoeconomicsUrl(chainId);
  const data = await fetch(
    `${ceUrl}/lm/?key=userData&address=${address.value}&timestamp=now`,
  );
  if (data.status !== 200) return null;
  const parsedData = await data.json();
  // pastRewards = dispensed
  // nextRewardPayment = claimed - dispensed
  // unclaimedReward = claimableReward - claimed
  // const parsedData = {
  //   timestamp: 200600,
  //   rewardBuckets: [],
  //   totalTickets: 56578387.68990869,
  //   user: {
  //     tickets: [
  //       {
  //         amount: 83.87962924761631,
  //         mul: 0.7274305555555691,
  //         reward: 21.384575625320533,
  //         timestamp: "April 15th 2021, 9:59:14 am",
  //       },
  //     ],
  //     claimed: 66.48991988613547,
  //     dispensed: 0,
  //     forfeited: 61.81223708007379,
  //     claimableReward: 82.04571361358246,
  //     reservedReward: 87.874495511456,
  //     totalTickets: 83.87962924761631,
  //     nextRewardShare: 0.0000014825383449832216,
  //     totalRewardAtMaturity: 2304.874495511456,
  //     ticketAmountAtMaturity: 83.87962924761631,
  //     yieldAtMaturity: 1.0476261793199715,
  //     maturityDate: "August 13th 2021, 9:59:14 am",
  //   },
  // };

  if (!parsedData.user.claimableReward) return null;
  return parsedData.user;
}

async function getVSData(address: ComputedRef<any>, chainId: string) {}
// `${ceUrl}/lm/?key=userData&address=${address.value}&timestamp=now`,

export default defineComponent({
  components: {
    Layout,
    SifButton,
    AssetItem,
    ActionsPanel,
    Copy,
    SubHeading,
    Box,
    Modal,
    ModalView,
    PairTable,
    Tooltip,
    Icon,
    ConfirmationModal,
  },
  methods: {
    openClaimModal() {
      this.modalOpen = true;
    },
    requestClose() {
      this.modalOpen = false;
    },
    claimRewards() {
      alert("claim logic/keplr goes here");
    },
  },
  data() {
    return {
      modalOpen: false,
      loadingLm: true,
    };
  },
  setup() {
    const { store, actions, config } = useCore();
    const address = computed(() => store.wallet.sif.address);
    const transactionState = ref<ConfirmState | string>("confirming");
    const transactionStateMsg = ref<string>("");
    const transactionHash = ref<string | null>(null);

    let lmRewards = ref<any>();
    let loadingLm = ref<Boolean>(true);

    watch(address, async () => {
      loadingLm.value = true;
      lmRewards.value = await getLMData(address, config.sifChainId);
      loadingLm.value = false;
    });

    onMounted(async () => {
      loadingLm.value = true;
      lmRewards.value = await getLMData(address, config.sifChainId);
      loadingLm.value = false;
    });

    async function handleAskConfirmClicked() {
      transactionState.value = "signing";
      const tx = await actions.clp.claimRewards();
      // transactionHash.value = tx.hash;
      // transactionState.value = toConfirmState(tx.state); // TODO: align states
      // transactionStateMsg.value = tx.memo ?? "";
    }

    const computedLMPairPanel = computed(() => {
      if (!lmRewards.value) {
        return [];
      }
      return [
        {
          key: "Claimable  Rewards",
          value: lmRewards.value.claimableReward,
        },
        {
          key: "Projected Full Amount",
          value: lmRewards.value.totalRewardAtMaturity,
        },
      ];
    });

    return {
      lmRewards,
      REWARD_INFO,
      computedLMPairPanel,
      format,
      handleAskConfirmClicked,
      transactionState,
      transactionStateMsg,
      transactionHash,
      loadingLm,
      address,
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
      <div v-if="loadingLm" class="loader-container">
        <div class="loader" />
      </div>

      <!-- TODO make this a component that can also handle VS-->
      <Box v-if="lmRewards">
        <div class="reward-container">
          <SubHeading>{{ REWARD_INFO["lm"].label }}</SubHeading>
          <Copy>
            {{ REWARD_INFO["lm"].description }}
          </Copy>
          <div class="details-container">
            <div class="amount-container">
              <div class="reward-rows">
                <div class="reward-row">
                  <div class="row-label">Claimable Rewards</div>
                  <div class="row-amount">
                    {{ format(lmRewards.claimableReward - lmRewards.claimed) }}
                  </div>
                  <AssetItem symbol="Rowan" :label="false" />
                </div>
                <div class="reward-row">
                  <div class="row-label">
                    Projected Full Amount
                    <Tooltip>
                      <template #message>
                        <div class="tooltip">
                          This is your projected Amount if all else equal until
                          Maturity Date <br />
                          Maturity Date <br />
                          {{ lmRewards.maturityDate }}
                        </div>
                      </template>
                      <Icon icon="info-box-black" />
                    </Tooltip>
                  </div>
                  <div class="row-amount">
                    {{ format(lmRewards.totalRewardAtMaturity) }}
                  </div>
                  <AssetItem symbol="Rowan" :label="false" />
                </div>
              </div>
            </div>
          </div>
          <div class="reward-buttons">
            <a
              class="more-info-button mr-8"
              target="_blank"
              :href="`https://cryptoeconomics.vercel.app/#${address}&type=lm`"
              >More Info</a
            >
            <SifButton @click="openClaimModal" :primary="true">Claim</SifButton>
          </div>
        </div>
      </Box>

      <Box v-else-if="!loadingLm && !lmRewards">No LM Rewards</Box>
    </div>
    <ActionsPanel connectType="connectToSif" />

    <div v-if="modalOpen">
      <ConfirmationModal
        :requestClose="requestClose"
        @confirmed="handleAskConfirmClicked"
        :state="transactionState"
        :transactionHash="transactionHash"
        :transactionStateMsg="transactionStateMsg"
        confirmButtonText="Claim Rewards"
        title="Claim Rewards"
      >
        <template v-slot:selecting>
          <div>
            <div class="claim-container">
              <Copy>
                Are you sure you want to claim your rewards? Once you claim
                these rewards, your multiplier will reset to 1x for all
                remaining amounts and will continue to accumulate if within the
                reward eligibility timeframe. 
                <br />
                <br />
                Please note that the rewards will be released at the end of the
                week.
                <br />
                <br />
                Find out <a href="">additional information here</a>.
              </Copy>
              <br />
              <PairTable :items="computedLMPairPanel" />
              <br />
              <!-- <div class="reward-buttons">
                <SifButton
                  class="reward-button"
                  @click="requestClose"
                  secondary="true"
                  >Cancel</SifButton
                >
                <SifButton
                  class="reward-button"
                  @click="claimRewards"
                  primary="true"
                  >Claim Rewards</SifButton
                >
              </div> -->
            </div>
          </div>
        </template>

        <template v-slot:common>
          <p class="text--normal" data-handle="confirmation-wait-message">
            Supplying
            <span class="text--bold">{{ fromAmount }} {{ fromSymbol }}</span>
            and
            <span class="text--bold">{{ toAmount }} {{ toSymbol }}</span>
          </p>
        </template>
      </ConfirmationModal>
    </div>
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
    }
    .details-container {
    }
  }
  .reward-rows {
    display: flex;
    flex-direction: column;
    margin-bottom: 15px;
  }
  .reward-row {
    display: flex;
    width: 100%;
    justify-content: space-between;

    font-size: 14px;
    font-weight: 400;
    .row-label {
      flex: 1 1 auto;
      text-align: left;
    }
    .row-amount {
      width: 100px;
      text-align: right;
    }
    .row {
      width: 15px;

      margin-left: 2px;
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

.reward-buttons {
  display: flex;
  flex-direction: row;

  justify-content: space-between;

  .more-info-button,
  .btn {
    width: 300px;
    display: block;
    font-weight: 600;
    font-style: italic;
  }
  .reward-button {
    text-align: center;
  }
}

/* MODAL Styles */
.claim-container {
  font-weight: 400;
  display: flex;
  flex-direction: column;
  padding: 30px 20px 20px 20px;
  min-height: 50vh;
  .container {
    font-size: 14px;
    line-height: 16px;
  }
}
</style>
