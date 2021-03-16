<script lang="ts">
import { defineComponent, ref, watch } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import {
  Asset,
  LiquidityProvider,
  PoolState,
  useRemoveLiquidityCalculator,
} from "ui-core";
import { useCore } from "@/hooks/useCore";
import { useRoute, useRouter } from "vue-router";
import { computed, effect, Ref, toRef } from "@vue/reactivity";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import Slider from "@/components/shared/Slider.vue";
import { toConfirmState } from "./utils/toConfirmState";
import { ConfirmState } from "@/types";
import ConfirmationModal from "@/components/shared/ConfirmationModal.vue";
import DetailsPanelRemove from "@/components/shared/DetailsPanelRemove.vue";

export default defineComponent({
  components: {
    AssetItem,
    Layout,
    ActionsPanel,
    Slider,
    ConfirmationModal,
    DetailsPanelRemove,
  },
  setup() {
    const { store, actions, poolFinder, api } = useCore();
    const route = useRoute();
    const router = useRouter();
    const transactionState = ref<ConfirmState>("selecting");
    const transactionHash = ref<string | null>(null);
    const transactionStateMsg = ref<string>("");
    const asymmetry = ref("0");
    const wBasisPoints = ref("0");
    const nativeAssetSymbol = ref("rowan");
    const externalAssetSymbol = ref<string | null>(
      route.params.externalAsset ? route.params.externalAsset.toString() : null,
    );
    const { connected } = useWalletButton();

    const liquidityProvider = ref(null) as Ref<LiquidityProvider | null>;
    const withdrawExternalAssetAmount: Ref<string | null> = ref(null);
    const withdrawNativeAssetAmount: Ref<string | null> = ref(null);
    const state = ref(0);

    effect(() => {
      if (!externalAssetSymbol.value) return null;
      api.ClpService.getLiquidityProvider({
        symbol: externalAssetSymbol.value,
        lpAddress: store.wallet.sif.address,
      }).then(liquidityProviderResult => {
        liquidityProvider.value = liquidityProviderResult;
      });
    });

    // if these values change, recalculate state and asset amounts
    watch([wBasisPoints, asymmetry, liquidityProvider], () => {
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
          !asymmetry.value
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
          return;

        transactionState.value = "signing";
        const tx = await actions.clp.removeLiquidity(
          Asset.get(externalAssetSymbol.value),
          wBasisPoints.value,
          asymmetry.value,
        );
        transactionHash.value = tx.hash;
        transactionState.value = toConfirmState(tx.state); // TODO: align states
        transactionStateMsg.value = tx.memo ?? "";
      },

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
      externalAssetSymbol,
      transactionState,
      transactionHash,
    };
  },
});
</script>

<template>
  <Layout
    class="pool"
    :backLink="`/pool/${externalAssetSymbol}`"
    title="Remove Liquidity"
  >
    <div :class="!withdrawNativeAssetAmount ? 'disabled' : 'active'">
      <div class="panel-header text--left">
        <div class="mb-10">Amount to Withdraw:</div>
        <h1>{{ wBasisPoints / 100 }}%</h1>
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
        <h4 class="text--left">You Should Receive:</h4>
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

    <ConfirmationModal
      :requestClose="requestTransactionModalClose"
      @confirmed="handleAskConfirmClicked"
      :state="transactionState"
      :transactionHash="transactionHash"
      :transactionStateMsg="transactionStateMsg"
      confirmButtonText="Confirm Withdrawal"
      title="You are withdrawing"
    >
      <template v-slot:selecting>
        <div>
          <DetailsPanelRemove
            class="details"
            :externalAssetSymbol="externalAssetSymbol"
            :externalAssetAmount="withdrawExternalAssetAmount"
            :nativeAssetSymbol="nativeAssetSymbol"
            :nativeAssetAmount="withdrawNativeAssetAmount"
          />
        </div>
      </template>

      <template v-slot:common>
        <p class="text--normal">
          You should receive
          <span class="text--bold">
            {{ withdrawExternalAssetAmount }}
            {{
              externalAssetSymbol.toLowerCase().includes("rowan")
                ? externalAssetSymbol.toUpperCase()
                : "c" + externalAssetSymbol.slice(1).toUpperCase()
            }}
          </span>
          and
          <span class="text--bold">
            {{ withdrawNativeAssetAmount }}
            {{ nativeAssetSymbol.toUpperCase() }}
          </span>
        </p>
      </template>
    </ConfirmationModal>
  </Layout>
</template>

<style lang="scss" scoped>
h1 {
  font-size: 42px;
  color: $c_gray_900;
}
.disabled {
  opacity: 0.3;
}
.panel-header {
  margin-bottom: 1.5rem;
}
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
