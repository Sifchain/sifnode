<script lang="ts">
import { defineComponent, ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import Layout from "@/components/layout/Layout.vue";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialogSif from "@/components/tokenSelector/SelectTokenDialogSif.vue";
import Modal from "@/components/shared/Modal.vue";
import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationDialog from "@/components/confirmationDialog/PoolConfirmationDialog.vue";
import { PoolState, usePoolCalculator } from "ui-core";
import { useCore } from "@/hooks/useCore";
import { useWallet } from "@/hooks/useWallet";
import { computed } from "@vue/reactivity";
import FatInfoTable from "@/components/shared/FatInfoTable.vue";
import FatInfoTableCell from "@/components/shared/FatInfoTableCell.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import { useCurrencyFieldState } from "@/hooks/useCurrencyFieldState";
import { toConfirmState } from "./utils/toConfirmState";
import { ConfirmState } from "../types";

export default defineComponent({
  components: {
    ActionsPanel,
    Layout,
    Modal,
    CurrencyPairPanel,
    SelectTokenDialogSif,
    FatInfoTable,
    ConfirmationDialog,
    FatInfoTableCell,
    ModalView,
  },
  props: ["title"],
  setup(props) {
    const { actions, poolFinder, store } = useCore();
    const selectedField = ref<"from" | "to" | null>(null);
    const transactionState = ref<ConfirmState>("selecting");
    const transactionHash = ref<string | null>(null);
    const router = useRouter();
    const route = useRoute();

    const { fromSymbol, fromAmount, toAmount } = useCurrencyFieldState();

    const toSymbol = ref("rowan");
    fromSymbol.value = route.params.externalAsset
      ? route.params.externalAsset.toString()
      : null;

    function clearAmounts() {
      fromAmount.value = "0.0";
      toAmount.value = "0.0";
    }

    const { connected, connectedText } = useWalletButton({
      addrLen: 8,
    });

    const { balances } = useWallet(store);

    const liquidityProvider = computed(() => {
      if (!fromSymbol) return null;
      return (
        store.accountpools.find((pool) => {
          pool.lp.asset.symbol === fromSymbol.value;
        })?.lp ?? null
      );
    });

    const {
      aPerBRatioMessage,
      bPerARatioMessage,
      aPerBRatioProjectedMessage,
      bPerARatioProjectedMessage,
      shareOfPoolPercent,
      totalLiquidityProviderUnits,
      tokenAFieldAmount,
      tokenBFieldAmount,
      preExistingPool,
      state,
    } = usePoolCalculator({
      balances,
      tokenAAmount: fromAmount,
      tokenBAmount: toAmount,
      tokenASymbol: fromSymbol,
      tokenBSymbol: toSymbol,
      poolFinder,
      liquidityProvider,
    });

    function handleNextStepClicked() {
      if (!tokenAFieldAmount.value)
        throw new Error("from field amount is not defined");
      if (!tokenBFieldAmount.value)
        throw new Error("to field amount is not defined");

      transactionState.value = "confirming";
    }

    async function handleAskConfirmClicked() {
      if (!tokenAFieldAmount.value)
        throw new Error("Token A field amount is not defined");
      if (!tokenBFieldAmount.value)
        throw new Error("Token B field amount is not defined");

      transactionState.value = "signing";
      const tx = await actions.clp.addLiquidity(
        tokenBFieldAmount.value,
        tokenAFieldAmount.value
      );

      transactionHash.value = tx.hash;
      transactionState.value = toConfirmState(tx.state); // TODO: align states

      clearAmounts();
    }

    function requestTransactionModalClose() {
      if (transactionState.value === "confirmed") {
        router.push("/pool");
      } else {
        transactionState.value = "selecting";
      }
    }

    return {
      fromAmount,
      fromSymbol,

      toAmount,
      toSymbol,

      connected,
      aPerBRatioMessage,
      bPerARatioMessage,
      aPerBRatioProjectedMessage,
      bPerARatioProjectedMessage,
      nextStepMessage: computed(() => {
        switch (state.value) {
          case PoolState.SELECT_TOKENS:
            return "Select Tokens";
          case PoolState.ZERO_AMOUNTS:
            return "Please enter an amount";
          case PoolState.INSUFFICIENT_FUNDS:
            return "Insufficient Funds";
          case PoolState.VALID_INPUT:
            return preExistingPool.value ? "Add liquidity" : "Create Pool";
        }
      }),
      nextStepAllowed: computed(() => {
        return state.value === PoolState.VALID_INPUT;
      }),
      handleFromSymbolClicked(next: () => void) {
        selectedField.value = "from";
        next();
      },
      handleToSymbolClicked(next: () => void) {
        selectedField.value = "to";
        next();
      },
      handleSelectClosed(data: string) {
        if (typeof data !== "string") {
          return;
        }

        if (selectedField.value === "from") {
          fromSymbol.value = data;
        }

        if (selectedField.value === "to") {
          toSymbol.value = data;
        }
        selectedField.value = null;
      },

      handleNextStepClicked,

      handleAskConfirmClicked,

      transactionHash,

      requestTransactionModalClose,

      transactionState,

      transactionModalOpen: computed(() => {
        return [
          "confirming",
          "signing",
          "failed",
          "rejected",
          "confirmed",
        ].includes(transactionState.value);
      }),

      handleBlur() {
        selectedField.value = null;
      },
      handleFromFocused() {
        selectedField.value = "from";
      },
      handleToFocused() {
        selectedField.value = "to";
      },
      handleFromMaxClicked() {
        selectedField.value = "from";
        const accountBalance = balances.value.find(
          (balance) => balance.asset.symbol === fromSymbol.value
        );
        if (!accountBalance) return;
        fromAmount.value = accountBalance.toFixed(8);
      },
      handleToMaxClicked() {
        selectedField.value = "to";
        const accountBalance = balances.value.find(
          (balance) => balance.asset.symbol === toSymbol.value
        );
        if (!accountBalance) return;
        toAmount.value = accountBalance.toFixed(8);
      },
      shareOfPoolPercent,
      connectedText,
      poolUnits: totalLiquidityProviderUnits,
    };
  },
});
</script>

<template>
  <Layout class="pool" backLink="/pool" :title="title">
    <Modal @close="handleSelectClosed">
      <template v-slot:activator="{ requestOpen }">
        <CurrencyPairPanel
          v-model:fromAmount="fromAmount"
          v-model:fromSymbol="fromSymbol"
          @fromfocus="handleFromFocused"
          @fromblur="handleBlur"
          @fromsymbolclicked="handleFromSymbolClicked(requestOpen)"
          :fromSymbolSelectable="connected"
          :fromMax="true"
          @frommaxclicked="handleFromMaxClicked"
          v-model:toAmount="toAmount"
          v-model:toSymbol="toSymbol"
          @tofocus="handleToFocused"
          @toblur="handleBlur"
          :toMax="true"
          @tomaxclicked="handleToMaxClicked"
          toSymbolFixed
          canSwapIcon="plus"
      /></template>
      <template v-slot:default="{ requestClose }">
        <SelectTokenDialogSif
          :selectedTokens="[fromSymbol, toSymbol].filter(Boolean)"
          @tokenselected="requestClose"
        />
      </template>
    </Modal>

    <FatInfoTable :show="nextStepAllowed">
      <template #header>Pool Token Prices</template>
      <template #body>
        <FatInfoTableCell>
          <span class="number">{{ aPerBRatioMessage }}</span
          ><br />
          <span
            >{{ fromSymbol.toUpperCase() }} per
            {{ toSymbol.toUpperCase() }}</span
          >
        </FatInfoTableCell>
        <FatInfoTableCell>
          <span class="number">{{ bPerARatioMessage }}</span
          ><br />
          <span
            >{{ toSymbol.toUpperCase() }} per
            {{ fromSymbol.toUpperCase() }}</span
          >
        </FatInfoTableCell>
      </template>
    </FatInfoTable>

    <FatInfoTable :show="nextStepAllowed">
      <template #header>Price Impact and Pool Share</template>
      <template #body>
        <FatInfoTableCell>
          <span class="number">{{ aPerBRatioProjectedMessage }}</span
          ><br />
          <span
            >{{ fromSymbol.toUpperCase() }} per
            {{ toSymbol.toUpperCase() }}</span
          >
        </FatInfoTableCell>
        <FatInfoTableCell>
          <span class="number">{{ bPerARatioProjectedMessage }}</span
          ><br />
          <span
            >{{ toSymbol.toUpperCase() }} per
            {{ fromSymbol.toUpperCase() }}</span
          >
        </FatInfoTableCell>
        <FatInfoTableCell>
          <span class="number">{{ shareOfPoolPercent }}</span
          ><br />Share of Pool
        </FatInfoTableCell></template
      >
    </FatInfoTable>

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
        :requestClose="requestTransactionModalClose"
        :fromToken="fromSymbol"
        :fromAmount="fromAmount"
        :poolUnits="poolUnits"
        :toAmount="toAmount"
        :toToken="toSymbol"
        :aPerB="aPerBRatioMessage"
        :bPerA="bPerARatioMessage"
        :shareOfPool="shareOfPoolPercent"
        :transactionHash="transactionHash"
    /></ModalView>
  </Layout>
</template>

<style lang="scss" scoped>
.number {
  font-size: 16px;
  font-weight: bold;
}
</style>
