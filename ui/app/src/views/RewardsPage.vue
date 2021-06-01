<script lang="ts">
import { computed, defineComponent, watch, onMounted } from "vue";
import { ref, ComputedRef } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { getCryptoeconomicsUrl } from "@/components/shared/utils";
import Layout from "@/components/layout/Layout.vue";
import ConfirmationModal from "@/components/shared/ConfirmationModal.vue";
import { Copy } from "@/components/shared/Text";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import Modal from "@/components/shared/Modal.vue";
import ModalView from "@/components/shared/ModalView.vue";
import PairTable from "@/components/shared/PairTable.vue";
import { ConfirmState } from "@/types";
import RewardContainer from "@/components/shared/RewardContainer/RewardContainer.vue";

async function getLMData(address: ComputedRef<any>, chainId: string) {
  if (!address.value) return;
  const ceUrl = getCryptoeconomicsUrl(chainId);
  const data = await fetch(
    `${ceUrl}/lm/?key=userData&address=${address.value}&timestamp=now`,
  );
  if (data.status !== 200) return {};
  const parsedData = await data.json();
  if (!parsedData.user || !parsedData.user) {
    return {};
  }
  return parsedData.user;
}

async function getVSData(address: ComputedRef<any>, chainId: string) {
  if (!address.value) return;
  const ceUrl = getCryptoeconomicsUrl(chainId);
  const data = await fetch(
    `${ceUrl}/vs/?key=userData&address=${address.value}&timestamp=now`,
  );
  if (data.status !== 200) return {};
  const parsedData = await data.json();
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
      this.modalOpen = true;
    },
    requestClose() {
      this.modalOpen = false;
    },
    claimRewards() {
      alert("claim logic/keplr goes here");
    },
    handleOpenModal(type: string) {
      console.log("type", type);
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
    const transactionHash = ref<string | null>(null);

    // TODO - We can do this better later
    let lmRewards = ref<any>();
    let vsRewards = ref<any>();
    let loadingVs = ref<Boolean>(true);

    watch(address, async () => {
      lmRewards.value = await getLMData(address, config.sifChainId);
      vsRewards.value = await getVSData(address, config.sifChainId);
    });

    onMounted(async () => {
      lmRewards.value = await getLMData(address, config.sifChainId);
      vsRewards.value = await getVSData(address, config.sifChainId);
    });

    async function handleAskConfirmClicked() {
      transactionState.value = "signing";
      // const tx = await actions.clp.claimRewards();
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
      computedLMPairPanel,
      computedVSPairPanel,
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
      <RewardContainer
        type="lm"
        :data="lmRewards"
        :address="address"
        @openModal="handleOpenModal"
      />
      <RewardContainer
        type="vs"
        :data="vsRewards"
        :address="address"
        @openModal="handleOpenModal"
      />
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
  padding: 30px 20px 20px 20px;
  min-height: 50vh;
  .container {
    font-size: 14px;
    line-height: 16px;
  }
}
</style>
