<script lang="ts">
import { computed, defineComponent, watch, onMounted } from "vue";
import { ref, ComputedRef } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { getCryptoeconomicsUrl, getLMData } from "@/components/shared/utils";
import Layout from "@/components/layout/Layout.vue";
import ConfirmationModal from "@/components/shared/ConfirmationModal.vue";
import { Copy } from "@/components/shared/Text";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import Modal from "@/components/shared/Modal.vue";
import ModalView from "@/components/shared/ModalView.vue";
import PairTable from "@/components/shared/PairTable.vue";
import { ConfirmState } from "@/types";
import RewardContainer from "@/components/shared/RewardContainer/RewardContainer.vue";
import { toConfirmState } from "./utils/toConfirmState";

const claimTypeMap = {
  lm: "2",
  vs: "3",
};
type IClaimType = "lm" | "vs" | null;

async function getVSData(address: ComputedRef<any>, chainId: string) {
  if (!address.value) return;
  // const ceUrl = getCryptoeconomicsUrl(chainId);
  // const data = await fetch(
  //   `${ceUrl}/vs/?key=userData&address=${address.value}&timestamp=now`,
  // );
  // if (data.status !== 200) return {};
  const parsedData = {
    totalDepositedAmount: 24205454.32624847,
    timestamp: 147600,
    rewardBuckets: [
      {
        rowan: 17301918.265222006,
        initialRowan: 45000000,
        duration: 1200,
      },
    ],
    user: {
      tickets: [
        {
          commission: 0,
          amount: 515.075,
          mul: 0.5112847222222296,
          reward: 259.84093497171625,
          validatorRewardAddress: "sif1avau6q23mrmmuz2nlyqk0mgrdzna5tf6yvc0et",
          validatorStakeAddress:
            "sifvaloper1avau6q23mrmmuz2nlyqk0mgrdzna5tf6dws9em",
          timestamp: "April 20th 2021, 11:28:43 pm",
          rewardDelta: 0.7991189486285039,
          poolDominanceRatio: 0.000021292080431235025,
          commissionRewardsByValidator: {
            sif1avau6q23mrmmuz2nlyqk0mgrdzna5tf6yvc0et: 0,
          },
        },
        {
          commission: 0,
          amount: 590,
          mul: 0.5052083333333406,
          reward: 289.3530824379835,
          validatorRewardAddress: "sif1lnhxf6war6qlldemkqzp0t3g57hpe9a664epyu",
          validatorStakeAddress:
            "sifvaloper1lnhxf6war6qlldemkqzp0t3g57hpe9a6nh3tyv",
          timestamp: "April 21st 2021, 10:48:43 pm",
          rewardDelta: 0.9153621893720666,
          poolDominanceRatio: 0.000024389317001269062,
          commissionRewardsByValidator: {
            sif1lnhxf6war6qlldemkqzp0t3g57hpe9a664epyu: 0,
          },
        },
      ],
      claimableRewardsOnWithdrawnAssets: 0,
      dispensed: 0,
      forfeited: 0,
      totalAccruedCommissionsAndClaimableRewards: 279.0362887823368,
      totalClaimableCommissionsAndClaimableRewards: 279.0362887823368,
      reservedReward: 549.1940174096998,
      totalDepositedAmount: 1105.075,
      totalClaimableRewardsOnDepositedAssets: 279.0362887823368,
      currentTotalCommissionsOnClaimableDelegatorRewards: 0,
      totalAccruedCommissionsAtMaturity: 0,
      totalCommissionsAndRewardsAtMaturity: 1340.756646783233,
      claimableCommissions: 0,
      delegatorAddresses: [],
      totalRewardsOnDepositedAssetsAtMaturity: 1340.756646783233,
      ticketAmountAtMaturity: 1105.075,
      yieldAtMaturity: 1.2132720826941457,
      nextRewardShare: 0.000045653966461668654,
      currentYieldOnTickets: 0.9607676926913524,
      maturityDate: "August 19th 2021, 10:48:43 pm",
      maturityDateISO: "2021-08-19T22:48:43.000Z",
      yearsToMaturity: 0.2168949771689498,
      currentAPYOnTickets: 4.4296447305138145,
      maturityDateMs: 0,
      futureReward: 1061.7203580008963,
      nextReward: 1.7120237423125746,
      nextRewardProjectedFutureReward: 4499.198394797446,
      nextRewardProjectedAPYOnTickets: 4.07139641634952,
      maturityAPY: 0,
    },
  };
  if (!parsedData.user || !parsedData.user) {
    return {};
  }
  return parsedData.user;
}

export default defineComponent({
  components: {
    Layout,
    ActionsPanel,
    Copy,
    Modal,
    ModalView,
    PairTable,
    ConfirmationModal,
    RewardContainer,
  },
  methods: {
    openClaimModal() {
      this.transactionState = "confirming";
    },
    requestClose() {
      this.transactionState = "selecting";
    },
    handleOpenModal(type: IClaimType) {
      console.log("type", type);
      this.claimType = type;
      this.openClaimModal();
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
    const { store, usecases, config } = useCore();
    const address = computed(() => store.wallet.sif.address);
    const transactionState = ref<ConfirmState | string>("selecting");
    const transactionStateMsg = ref<string>("");
    const transactionHash = ref<string | null>(null);
    // TODO - We can do this better later
    let lmRewards = ref<any>();
    let vsRewards = ref<any>();
    let loadingVs = ref<Boolean>(true);
    let claimType = ref<IClaimType>(null);

    watch(address, async () => {
      lmRewards.value = await getLMData(address, config.sifChainId);
      vsRewards.value = await getVSData(address, config.sifChainId);
    });

    onMounted(async () => {
      lmRewards.value = await getLMData(address, config.sifChainId);
      vsRewards.value = await getVSData(address, config.sifChainId);
    });

    async function handleAskConfirmClicked() {
      if (!claimType.value) {
        return console.error("No claim type");
      }
      transactionState.value = "signing";
      const tx = await usecases.dispensation.claim({
        fromAddress: address.value,
        claimType: claimTypeMap[claimType.value] as "2" | "3",
      });
      transactionHash.value = tx.hash;
      transactionState.value = toConfirmState(tx.state); // TODO: align states
      transactionStateMsg.value = tx.memo ?? "";
    }

    const computedPairPanel = computed(() => {
      if (!claimType.value) {
        return console.error("No claim type");
      }
      let data;
      claimType.value === "lm" ? (data = lmRewards) : (data = vsRewards);
      return [
        {
          key: "Claimable  Rewards",
          value: data.value.totalClaimableCommissionsAndClaimableRewards,
        },
        {
          key: "Projected Full Amount",
          value: data.value.totalCommissionsAndRewardsAtMaturity,
        },
        {
          key: "Maturity Date",
          value: data.value.maturityDateISO,
          type: "date",
        },
      ];
    });

    return {
      lmRewards,
      vsRewards,
      computedPairPanel,
      handleAskConfirmClicked,
      transactionState,
      transactionStateMsg,
      transactionHash,
      loadingVs,
      address,
      claimType,
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
      <RewardContainer
        claimType="lm"
        :data="lmRewards"
        :address="address"
        :claimDisabled="false"
        @openModal="handleOpenModal"
      />
      <RewardContainer
        claimType="vs"
        :data="vsRewards"
        :address="address"
        :claimDisabled="false"
        @openModal="handleOpenModal"
      />
    </div>

    <ActionsPanel connectType="connectToSif" />

    <div v-if="transactionState !== 'selecting'">
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
          <div class="claim-container">
            <Copy class="mb-8">
              Are you sure you want to claim your rewards? Claiming your rewards
              will reset all of your tickets at this very moment. Resetting your
              tickers will release your rewards based on its current multiplier.
              Reset tickets then start empty with a 25% multiplier again and
              will continue to accumulate if within the reward eligibility
              timeframe. Unless you have reached full maturity, we recommend
              that you do not claim so you can realize your full rewards.
            </Copy>
            <Copy class="mb-8">
              Please note that the rewards will be dispensed at the end of the
              week.
            </Copy>
            <Copy class="mb-8">
              Find out
              <a
                href="https://docs.sifchain.finance/resources/rewards-programs"
                target="_blank"
                >additional information here</a
              >.
            </Copy>
            <PairTable :items="computedPairPanel" class="mb-10" />
          </div>
        </template>

        <template v-slot:common>
          <p class="text--normal" data-handle="confirmation-wait-message">
            {{ claimType === "lm" ? "Liquidity Mining" : "Validator Subsidy" }}
            Rewards <br /><br />
            Claim {{ computedPairPanel[0].value }} Rowan
          </p>
        </template>
      </ConfirmationModal>
    </div>
  </Layout>
</template>

<style scoped lang="scss">
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
}

/* MODAL Styles */
.claim-container {
  font-weight: 400;
  display: flex;
  flex-direction: column;
  // padding: 30px 20px 20px 20px;
  min-height: 50vh;
  .container {
    font-size: 14px;
    line-height: 21px;
  }
}
</style>
