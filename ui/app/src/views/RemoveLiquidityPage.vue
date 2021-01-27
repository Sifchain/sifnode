<script lang="ts">
import { defineComponent, ref, watch } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialogSif from "@/components/tokenSelector/SelectTokenDialogSif.vue";
import Modal from "@/components/shared/Modal.vue";
import ModalView from "@/components/shared/ModalView.vue";
import { Asset, PoolState, useRemoveLiquidityCalculator } from "ui-core";
import { LiquidityProvider } from "ui-core";
import { useCore } from "@/hooks/useCore";
import { useRoute, useRouter } from "vue-router";
import { computed, effect, Ref, toRef } from "@vue/reactivity";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import SifButton from "@/components/shared/SifButton.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import Caret from "@/components/shared/Caret.vue";
import Slider from "@/components/shared/Slider.vue";
import ConfirmationDialog, {
  ConfirmState,
} from "@/components/confirmationDialog/RemoveConfirmationDialog.vue";

export default defineComponent({
  components: {
    AssetItem,
    Layout,
    Modal,
    ModalView,
    SelectTokenDialogSif,
    ActionsPanel,
    SifButton,
    Caret,
    Slider,
    ConfirmationDialog,
  },
  setup() {
    const { store, actions, poolFinder, api } = useCore();
    const route = useRoute();
    const router = useRouter();
    const transactionState = ref<ConfirmState>("selecting");
    const transactionHash = ref<string | null>(null);
    const asymmetry = ref("0");
    const wBasisPoints = ref("0");
    const nativeAssetSymbol = ref("rowan");
    const externalAssetSymbol = ref<string | null>(
      route.params.externalAsset ? route.params.externalAsset.toString() : null
    );
    const { connected, connectedText } = useWalletButton({
      addrLen: 8,
    });

    const liquidityProvider = ref(null) as Ref<LiquidityProvider | null>;
    let withdrawExternalAssetAmount: Ref<string | null> = ref(null)
    let withdrawNativeAssetAmount: Ref<string | null> = ref(null)
    let state = ref(0)

    effect(() => {
      if (!externalAssetSymbol.value) return null;
      api.ClpService.getLiquidityProvider({
        symbol: externalAssetSymbol.value,
        lpAddress: store.wallet.sif.address,
      }).then((liquidityProviderResult) => {
        liquidityProvider.value = liquidityProviderResult;
      });
    });

    // if these values change, recalculate state and asset amounts
    watch([
      wBasisPoints, 
      asymmetry, 
      liquidityProvider
    ], () => {
      const calcData = useRemoveLiquidityCalculator({
        externalAssetSymbol,
        nativeAssetSymbol,
        wBasisPoints,
        asymmetry,
        liquidityProvider,
        sifAddress: toRef(store.wallet.sif, "address"),
        poolFinder,
      });
      state.value = calcData.state;
      withdrawExternalAssetAmount.value = calcData.withdrawExternalAssetAmount;
      withdrawNativeAssetAmount.value = calcData.withdrawNativeAssetAmount;
    });
   
    return {
      connected,
      state,
      nextStepMessage: computed(() => {
        switch (state.value) {
          case PoolState.SELECT_TOKENS:
            return "Select Tokens";
          case PoolState.ZERO_AMOUNTS:
            return "Please enter an amount";
          case PoolState.NO_LIQUIDITY:
            return "No liquidity available.";
          case PoolState.INSUFFICIENT_FUNDS:
            return "Insufficient funds in this pool";
          case PoolState.VALID_INPUT:
            return "Remove Liquidity";
        }
      }),
      nextStepAllowed: computed(() => {
        return state.value === PoolState.VALID_INPUT;
      }),
      handleSelectClosed(data: string | MouseEvent) {
        if (typeof data !== "string") {
          return;
        }

        externalAssetSymbol.value = data;
      },
      handleNextStepClicked() {
        if (
          !externalAssetSymbol.value ||
          !wBasisPoints.value ||
          !asymmetry.value ||
          state.value !== PoolState.VALID_INPUT
        )
          return;

        transactionState.value = "confirming";
      },
      async handleAskConfirmClicked() {
        if (
          !externalAssetSymbol.value ||
          !wBasisPoints.value ||
          !asymmetry.value
        ) 
          return
          
        transactionState.value = "signing";
        try {
          let tx = await actions.clp.removeLiquidity(
            Asset.get(externalAssetSymbol.value),
            wBasisPoints.value,
            asymmetry.value
          );
          transactionHash.value = tx?.transactionHash ?? "";
          transactionState.value = "confirmed";
        } catch (err) {
          transactionState.value = "failed";
        }
      },

      transactionModalOpen: computed(() => {
        return ["confirming", "signing", "confirmed"].includes(
          transactionState.value
        );
      }),

      requestTransactionModalClose() {
        if (transactionState.value === "confirmed") {
          router.push("/pool");
        } else {
          transactionState.value = "selecting";
        }
      },
      PoolState,
      wBasisPoints,
      asymmetry,
      nativeAssetSymbol,
      withdrawExternalAssetAmount,
      withdrawNativeAssetAmount,
      connectedText,
      externalAssetSymbol,
      transactionState,
      transactionHash,
    };
  },
});
</script>

<template>
  <Layout class="pool" :backLink='`/pool/${externalAssetSymbol}`' title="Remove Liquidity"  >
  <div :class="!withdrawNativeAssetAmount ? 'disabled' : 'active' ">

    <div class="panel-header text--left">
      <div class="mb-10">Amount:</div>
      <h1>{{wBasisPoints/100}}%</h1>
    </div>

    <Slider
      message=""
      :disabled="!connected || state === PoolState.NO_LIQUIDITY"
      v-model="wBasisPoints"
      min="0"
      max="10000"
      type="range"
      step="1"
      @leftclicked="wBasisPoints = '0'"
      @middleclicked="wBasisPoints = '5000'"
      @rightclicked="wBasisPoints = '10000'"
      leftLabel="0%"
      middleLabel="50%"
      rightLabel="100%"
    />

    <Slider
      class="pt-4"
      message="Choose which ratio to withdraw from each asset"
      :disabled="!connected || state === PoolState.NO_LIQUIDITY"
      v-model="asymmetry"
      min="-10000"
      max="10000"
      type="range"
      step="1"
      @leftclicked="asymmetry = '-10000'"
      @middleclicked="asymmetry = '0'"
      @rightclicked="asymmetry = '10000'"
      leftLabel="All Rowan"
      middleLabel="Equal"
      rightLabel="All External Asset"
    />
    <div class="asset-row">
      <h4 class="text--left">Total Deposited After Transaction</h4>
      <div>
        <AssetItem :symbol="nativeAssetSymbol" />
        <AssetItem :symbol="externalAssetSymbol" />
      </div>
      <div>
        <div>{{ withdrawNativeAssetAmount }}</div>
        <div>{{ withdrawExternalAssetAmount }}</div>
      </div>
    </div>
    </div>

    <ActionsPanel
      @nextstepclick="handleNextStepClicked"
      :nextStepAllowed="nextStepAllowed"
      :nextStepMessage="nextStepMessage"
    />

    <ModalView
      :requestClose="requestTransactionModalClose"
      :isOpen="transactionModalOpen"
      ><ConfirmationDialog
        @confirmswap="handleAskConfirmClicked"
        :state="transactionState"
        :externalAssetSymbol="externalAssetSymbol"
        :nativeAssetSymbol="nativeAssetSymbol"
        :externalAssetAmount="withdrawExternalAssetAmount"
        :nativeAssetAmount="withdrawNativeAssetAmount"
        :transactionHash="transactionHash"
        :requestClose="requestTransactionModalClose"
    /></ModalView>
  </Layout>
</template>

<style lang="scss" scoped>
h1 { font-size: 42px; color: $c_gray_900}
.disabled {
    opacity: .3
  }
.panel-header {margin-bottom: 1.5rem}
.asset-row {
  margin-top: 1rem;
  margin-bottom: 1rem;
  background: $c_white;
  padding: 8px 8px 16px 8px;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  div {
    display: flex;
    justify-content: space-between;
  }
}
</style>
