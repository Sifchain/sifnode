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
import { format } from "ui-core/src/utils/format";
import Loader from "@/components/shared/Loader.vue";

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

async function getLMData(address: ComputedRef<any>, chainId: string) {
  if (!address.value) return;
  // const ceUrl = getCryptoeconomicsUrl(chainId);
  // const data = await fetch(
  //   `${ceUrl}/lm/?key=userData&address=${address.value}&timestamp=now`,
  // );
  // if (data.status !== 200) return {};
  // const parsedData = await data.json();
  // if (!parsedData.user || !parsedData.user.claimableReward) {
  //   return {};
  // }
  // return parsedData.user;
  return {
    claimableReward: 10000,
    claimed: 1133.5803574233153,
    dispensed: 0,
    forfeited: 2820.4719248996316,
    nextRewardProjectedAPYOnTickets: 1.56,
    maturityDate: "June 12th 2021, 11:19:14 am",
    nextRewardShare: 0,
    reservedReward: 1133.5803574233153,
    ticketAmountAtMaturity: 0,
    tickets: [],
    totalRewardAtMaturity: 1133.5803574233153,
    totalTickets: 0,
    yieldAtMaturity: null,
  };
}

async function getVSData(address: ComputedRef<any>, chainId: string) {
  if (!address.value) return;
  // const timestamp = Date.parse(new Date().toString());
  const ceUrl = getCryptoeconomicsUrl(chainId);
  const data = await fetch(
    `${ceUrl}/vs/?key=userData&address=${address.value}&timestamp=now`,
  );
  if (data.status !== 200) return null;
  const parsedData = await data.json();
  // TODO - VS endpoint does not return the same thing as LM endpoint so
  // mocking the data in the interim;
  return {
    claimableReward: 1133.5803574233153,
    claimed: 1133.5803574233153,
    dispensed: 0,
    forfeited: 2820.4719248996316,
    maturityDate: "June 12th 2021, 11:19:14 am",
    nextRewardShare: 0,
    reservedReward: 1133.5803574233153,
    ticketAmountAtMaturity: 0,
    tickets: [],
    totalRewardAtMaturity: 1133.5803574233153,
    totalTickets: 0,
    yieldAtMaturity: null,
  };
  // if (!parsedData.user.claimableReward) return null;
  // return parsedData.user;
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
    Modal,
    ModalView,
    PairTable,
    Tooltip,
    Icon,
    ConfirmationModal,
    Loader,
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
      loadingVs: true,
    };
  },
  setup() {
    const { store, actions, config } = useCore();
    const address = computed(() => store.wallet.sif.address);
    const transactionState = ref<ConfirmState | string>("confirming");
    const transactionStateMsg = ref<string>("");
    const transactionHash = ref<any>(null);

    // TODO - We can do this better later
    let lmRewards = ref<any>();
    let vsRewards = ref<any>();
    let loadingVs = ref<Boolean>(true);

    watch(address, async () => {
      lmRewards.value = await getLMData(address, config.sifChainId);
      // loadingVs.value = true;
      // vsRewards.value = await getVSData(address, config.sifChainId);
      // loadingVs.value = false;
    });

    onMounted(async () => {
      lmRewards.value = await getLMData(address, config.sifChainId);
      // loadingVs.value = true;
      // vsRewards.value = await getVSData(address, config.sifChainId);
      // loadingVs.value = false;
    });

    async function handleAskConfirmClicked() {
      transactionState.value = "signing";
      const tx = await actions.dispensation.claimRewards();
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

    const computedVSPairPanel = computed(() => {
      if (!vsRewards.value) {
        return [];
      }
      console.log("vsRewards", vsRewards);
      return [
        {
          key: "Claimable  Rewards",
          value: vsRewards.value.claimableReward,
        },
        {
          key: "Projected Full Amount",
          value: vsRewards.value.totalRewardAtMaturity,
        },
      ];
    });
    return {
      lmRewards,
      vsRewards,
      REWARD_INFO,
      computedLMPairPanel,
      computedVSPairPanel,
      format,
      handleAskConfirmClicked,
      transactionState,
      transactionStateMsg,
      transactionHash,
      loadingVs,
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
        class="underline"
        href="https://docs.sifchain.finance/resources/rewards-programs"
        >current rewards program</a
      >
      and how to become eligible.
    </Copy>
    <div class="rewards-container">
      <!-- TODO make this a component that can also handle VS && DRY ME -->
      <Box v-if="true">
        <div class="reward-container">
          <SubHeading>{{ REWARD_INFO["lm"].label }}</SubHeading>
          <Copy>
            {{ REWARD_INFO["lm"].description }}
          </Copy>
          <div class="details-container">
            <Loader v-if="!lmRewards" black />

            <div v-else class="amount-container">
              <div class="reward-rows">
                <div class="reward-row">
                  <div class="row-label">Claimable Rewards</div>
                  <div class="row-amount">
                    {{
                      format(lmRewards.claimableReward - lmRewards.claimed, {
                        mantissa: 4,
                      }) || "0"
                    }}
                  </div>
                  <AssetItem symbol="Rowan" :label="false" />
                </div>

                <div class="reward-row">
                  <div class="row-label">
                    Pending Reward Dispensation
                    <Tooltip>
                      <template #message>
                        <div class="tooltip">
                          Rewards that have been claimed and are pending
                          dispensation due to removal of liquidity or
                          user-initiated claims. Pending rewards are dispensed
                          every Friday once we have the dispensation module up
                          and running.
                        </div>
                      </template>
                      <Icon icon="info-box-black" />
                    </Tooltip>
                  </div>
                  <div class="row-amount">
                    {{
                      format(lmRewards.claimed - lmRewards.dispensed, {
                        mantissa: 4,
                      }) || "0"
                    }}
                  </div>
                  <AssetItem symbol="Rowan" :label="false" />
                </div>

                <div class="reward-row">
                  <div class="row-label">
                    Dispensed Rewards
                    <Tooltip>
                      <template #message>
                        <div class="tooltip">
                          Rewards that have already been dispensed.
                        </div>
                      </template>
                      <Icon icon="info-box-black" />
                    </Tooltip>
                  </div>
                  <div class="row-amount">
                    {{ format(lmRewards.dispensed, { mantissa: 4 }) || "0" }}
                  </div>
                  <AssetItem symbol="Rowan" :label="false" />
                </div>

                <div class="reward-row secondary">
                  <div class="row-label">
                    Projected Full Amount
                    <Tooltip>
                      <template #message>
                        <div class="tooltip">
                          <div v-if="lmRewards.maturityDate">
                            Projected Full Maturity Date: <br />
                            <span class="tooltip-date">{{
                              lmRewards.maturityDate
                            }}</span>
                            <span
                              v-if="lmRewards.nextRewardProjectedAPYOnTickets"
                            >
                              Projected Fully Maturated APY: <br />
                              <span class="tooltip-date">
                                {{
                                  format(
                                    lmRewards.nextRewardProjectedAPYOnTickets *
                                      100,
                                    {
                                      mantissa: 2,
                                    },
                                  )
                                }}%</span
                              >
                            </span>
                            <br /><br />
                          </div>
                          This is your estimated projected full reward amount
                          that you can earn if you were to leave your current
                          liquidity positions in place to the above mentioned
                          date. This includes projected future rewards, and
                          already claimed/disbursed previous rewards. This
                          number can fluctuate due to other market conditions
                          and this number is a representation of the current
                          market as it is in this very moment.
                        </div>
                      </template>
                      <Icon icon="info-box-black" />
                    </Tooltip>
                  </div>
                  <div class="row-amount">
                    {{
                      format(lmRewards.totalRewardAtMaturity, {
                        mantissa: 4,
                      }) || "0"
                    }}
                  </div>
                  <AssetItem symbol="Rowan" :label="false" />
                </div>
              </div>
              <div class="reward-buttons">
                <a
                  class="more-info-button mr-8"
                  target="_blank"
                  :href="`https://cryptoeconomics.sifchain.finance/#${address}&type=lm`"
                  >More Info</a
                >

                <!-- :disabled="(lmRewards.claimableReward - lmRewards.claimed) === 0" -->
                <SifButton
                  @click="openClaimModal"
                  :primary="true"
                  :disabled="false"
                  >Claim</SifButton
                >
              </div>
            </div>
          </div>
        </div>
      </Box>

      <!-- Validator Subsidy -->
      <Box v-if="vsRewards">
        <div class="reward-container">
          <SubHeading>Validator Subsidy</SubHeading>
          <Copy>
            Missing copy here. Lorem ipsum dolor sit amet, consectetur
            adipiscing elit, sed do eiusmod tempor incididunt ut labore et
            dolore magna aliqua. Ut enim ad minim veniam.
          </Copy>
          <div class="details-container">
            <div class="amount-container">
              <div class="reward-rows">
                <div class="reward-row">
                  <div class="row-label">Claimable Rewards</div>
                  <div class="row-amount">
                    {{
                      format(lmRewards.claimableReward - lmRewards.claimed, {
                        mantissa: 4,
                      })
                    }}
                  </div>
                  <AssetItem symbol="Rowan" :label="false" />
                </div>
                <div class="reward-row">
                  <div class="row-label">
                    Projected Full Amount
                    <Tooltip>
                      <template #message>
                        <div class="tooltip">
                          Projected Full Maturity Date: <br />
                          <span class="tooltip-date">{{
                            lmRewards.maturityDate
                          }}</span
                          ><br /><br />
                          This is your estimated projected full reward amount
                          that you can earn if you were to leave your current
                          liquidity positions in place to the above mentioned
                          date. This number can fluctuate due to other market
                          conditions and this number is a representation of the
                          current market as it is in this very moment.
                        </div>
                      </template>
                      <Icon icon="info-box-black" />
                    </Tooltip>
                  </div>
                  <div class="row-amount">
                    {{
                      format(lmRewards.totalRewardAtMaturity, { mantissa: 4 })
                    }}
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
                <p>
                  Are you sure you want to claim your rewards? Once you claim
                  these rewards, additional rewards will not accumulate.
                </p>
                <p>
                  Although if you re-add liquidity during our program window,
                  you will start accumulating rewards again from a 1x
                  multiplier.
                </p>
                <p>
                  Please note that the rewards will be released at the end of
                  the week.
                </p>
                <p>
                  Find out
                  <a href="" class="underline">additional information here</a>.
                </p>
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
    color: #343434;
  }
  .reward-row {
    display: flex;
    width: 100%;
    justify-content: space-between;
    font-size: $fs;
    font-weight: 400;
    &.secondary {
      color: #818181;
    }
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
  .more-info-button {
    background: #f3f3f3;
    color: #343434;
    font-weight: 100;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .more-info-button,
  .btn {
    width: 300px;
    border-radius: 6px;
    display: flex;
    font-size: $fs;
    height: 30px;
  }
  .reward-button {
    text-align: center;
  }
}

.tooltip-date {
  font-weight: 600;
}

/* MODAL Styles */
.claim-container {
  font-weight: 400;
  display: flex;
  flex-direction: column;

  // padding: 30px 20px 20px 20px;
  // min-height: 50vh;
  p {
    margin-bottom: 16px;
    font-size: 16px;
    line-height: 1.3;
    color: #343434;
  }
  .container {
    font-size: 14px;
    line-height: 16px;
  }
}
</style>
