<script lang="ts">
import { defineComponent, ref } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialog from "@/components/tokenSelector/SelectTokenDialog.vue";
import Modal from "@/components/shared/Modal.vue";
import {
  Asset,
  LiquidityProvider,
  PoolState,
  useRemoveLiquidityCalculator,
} from "../../../core";
import { useCore } from "@/hooks/useCore";

import { computed, toRef } from "@vue/reactivity";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import SifButton from "@/components/shared/SifButton.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import Caret from "@/components/shared/Caret.vue";
import { Fraction } from "../../../core/src/entities/fraction/Fraction";

export default defineComponent({
  components: {
    AssetItem,
    Layout,
    Modal,
    SelectTokenDialog,
    ActionsPanel,
    SifButton,
    Caret,
  },
  setup() {
    const { store, api } = useCore();
    const marketPairFinder = api.MarketService.find;
    const liquidityProviderFinder = (asset: Asset, address: string) =>
      LiquidityProvider(asset, new Fraction("10000"), address);

    const asymmetry = ref("0");
    const wBasisPoints = ref("5000");
    const nativeAssetSymbol = ref("rwn");
    const externalAssetSymbol = ref<string | null>(null);
    const { connected, connectedText } = useWalletButton({
      addrLen: 8,
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
      liquidityProviderFinder,
      sifAddress: toRef(store.wallet.sif, "address"),
      marketPairFinder,
    });
    // input not updating for some reason?

    return {
      connected,

      nextStepMessage: computed(() => {
        switch (state.value) {
          case PoolState.SELECT_TOKENS:
            return "Select Tokens";
          case PoolState.ZERO_AMOUNTS:
            return "Please enter an amount";
          case PoolState.INSUFFICIENT_FUNDS:
            return "Amount to remove is too high";
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
        alert(`Remove Liquidity!`);
      },
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
    <div class="slider">
      <p>Choose from 0 to 100% of how much to withdraw</p>
      <input
        v-model="wBasisPoints"
        class="input"
        type="range"
        max="10000"
        min="0"
        step="1"
      />
      <div class="row">
        <div @click="wBasisPoints = '0'">0%</div>
        <div @click="wBasisPoints = '5000'">50%</div>
        <div @click="wBasisPoints = '10000'">100%</div>
      </div>
    </div>
    <div class="slider">
      <p>Choose how much to withdraw from each asset</p>
      <input
        v-model="asymmetry"
        class="input"
        min="-10000"
        max="10000"
        type="range"
        step="1"
      />
      <div class="row">
        <div @click="asymmetry = '-10000'">All Rowan</div>
        <div @click="asymmetry = '0'">Equal</div>
        <div @click="asymmetry = '10000'">All Asset</div>
      </div>
    </div>
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
            <SelectTokenDialog
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
.slider {
  margin-bottom: 1rem;
  width: 100%;
  .input {
    width: 100%;
  }
  .row {
    display: flex;
    justify-content: space-between;
    & > * {
      width: 20%;
    }
    & > *:first-child {
      text-align: left;
    }
    & > *:last-child {
      text-align: right;
    }
  }
}
</style>
