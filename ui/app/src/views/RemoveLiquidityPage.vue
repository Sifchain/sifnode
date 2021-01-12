<script lang="ts">
import { defineComponent, ref } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialogSif from "@/components/tokenSelector/SelectTokenDialogSif.vue";
import Modal from "@/components/shared/Modal.vue";
import { Asset, PoolState, useRemoveLiquidityCalculator } from "ui-core";
import { LiquidityProvider } from "ui-core";
import { useCore } from "@/hooks/useCore";

import { computed, effect, Ref, toRef } from "@vue/reactivity";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import SifButton from "@/components/shared/SifButton.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import Caret from "@/components/shared/Caret.vue";
import { Fraction } from "ui-core";
import Slider from "@/components/shared/Slider.vue";

export default defineComponent({
  components: {
    AssetItem,
    Layout,
    Modal,
    SelectTokenDialogSif,
    ActionsPanel,
    SifButton,
    Caret,
    Slider,
  },
  setup() {
    const { store, actions, poolFinder, api } = useCore();

    const asymmetry = ref("0");
    const wBasisPoints = ref("5000");
    const nativeAssetSymbol = ref("rowan");
    const externalAssetSymbol = ref<string | null>(null);
    const { connected, connectedText } = useWalletButton({
      addrLen: 8,
    });

    const liquidityProvider = ref(null) as Ref<LiquidityProvider | null>;

    effect(() => {
      if (!externalAssetSymbol.value) return null;

      api.ClpService.getLiquidityProvider({
        symbol: externalAssetSymbol.value,
        lpAddress: store.wallet.sif.address,
      }).then((liquidityProviderResult) => {
        liquidityProvider.value = liquidityProviderResult;
      });
    });

    const {
      withdrawExternalAssetAmount,
      withdrawNativeAssetAmount,
      state,
    } = useRemoveLiquidityCalculator({
      externalAssetSymbol,
      nativeAssetSymbol,
      wBasisPoints,
      asymmetry,
      liquidityProvider,
      sifAddress: toRef(store.wallet.sif, "address"),
      poolFinder,
    });
    // input not updating for some reason?
    function clearFields() {
      asymmetry.value = "0";
      wBasisPoints.value = "0";
      nativeAssetSymbol.value = "rowan";
      externalAssetSymbol.value = null;
    }
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
      async handleNextStepClicked() {
        if (
          !externalAssetSymbol.value ||
          !wBasisPoints.value ||
          !asymmetry.value
        )
          return;

        try {
          await actions.clp.removeLiquidity(
            Asset.get(externalAssetSymbol.value),
            wBasisPoints.value,
            asymmetry.value
          );
          alert("Liquidity Removed");
        } catch (err) {
          alert(err);
        }
        clearFields();
      },
      PoolState,
      wBasisPoints,
      asymmetry,
      nativeAssetSymbol,
      withdrawExternalAssetAmount,
      withdrawNativeAssetAmount,
      connectedText,
      externalAssetSymbol,
    };
  },
});
</script>

<template>
  <Layout class="pool" backLink="/pool">
    <Slider
      message="Choose from 0 to 100% of how much to withdraw"
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
      message="Choose how much to withdraw from each asset"
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
      rightLabel="All Asset"
    />
    <div class="asset-row">
      <AssetItem :symbol="nativeAssetSymbol" />
      <div class="select-asset">
        <Modal @close="handleSelectClosed">
          <template v-slot:activator="{ requestOpen }">
            <SifButton
              v-if="externalAssetSymbol !== null"
              block
              @click="requestOpen"
            >
              <span><AssetItem :symbol="externalAssetSymbol" /></span>
              <span><Caret /></span>
            </SifButton>
            <SifButton
              v-if="externalAssetSymbol === null"
              primary
              block
              :disabled="!connected"
              @click="requestOpen"
            >
              <span>Select</span>
            </SifButton>
          </template>
          <template v-slot:default="{ requestClose }">
            <SelectTokenDialogSif
              :selectedTokens="[externalAssetSymbol].filter(Boolean)"
              @tokenselected="requestClose"
            />
          </template>
        </Modal>
      </div>
    </div>
    <div class="asset-row">
      <div>{{ withdrawNativeAssetAmount }}</div>
      <div>{{ withdrawExternalAssetAmount }}</div>
    </div>

    <ActionsPanel
      @nextstepclick="handleNextStepClicked"
      :nextStepAllowed="nextStepAllowed"
      :nextStepMessage="nextStepMessage"
    />
  </Layout>
</template>

<style lang="scss" scoped>
.asset-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 1rem;
}
</style>
