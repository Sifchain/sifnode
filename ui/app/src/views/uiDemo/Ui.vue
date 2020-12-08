<script lang="ts">
import { defineComponent, ref } from "vue";
import Layout from "@/components/layout/Layout.vue";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialog from "@/components/tokenSelector/SelectTokenDialog.vue";
import Modal from "@/components/shared/Modal.vue";
import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationDialog, {
  ConfirmState,
} from "@/components/confirmationDialog/PoolConfirmationDialog.vue";
import { PoolState, usePoolCalculator } from "ui-core";
import { useCore } from "@/hooks/useCore";
import { useWallet } from "@/hooks/useWallet";
import { computed } from "@vue/reactivity";
import PriceCalculation from "@/components/shared/PriceCalculation.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import { useCurrencyFieldState } from "@/hooks/useCurrencyFieldState";

export default defineComponent({
  components: {
    ActionsPanel,
    Layout,
    Modal,
    ModalView,
    ConfirmationDialog,
    CurrencyPairPanel,
    SelectTokenDialog,
    PriceCalculation,
  },
  setup() {
    const { actions, store, api } = useCore();
    const marketPairFinder = api.MarketService.find;
    const selectedField = ref<"from" | "to" | null>(null);

    const {
      fromSymbol,
      fromAmount,
      toSymbol,
      toAmount,
    } = useCurrencyFieldState();

    toSymbol.value = "rwn";

    const priceMessage = ref("");

    function clearAmounts() {
      fromAmount.value = "0.0";
      toAmount.value = "0.0";
    }

    const {
      connected,

      connectedText,
    } = useWalletButton({
      addrLen: 8,
    });

    const { balances } = useWallet(store);

    const {
      aPerBRatioMessage,
      bPerARatioMessage,
      shareOfPoolPercent,
      fromFieldAmount,
      toFieldAmount,
      preExistingPool,
      state,
    } = usePoolCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      marketPairFinder,
    });

    return {
      fromAmount,
      fromSymbol,

      toAmount,
      toSymbol,

      priceMessage,
      connected,
      aPerBRatioMessage,
      bPerARatioMessage,

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
      async handleNextStepClicked() {
        if (!fromFieldAmount.value)
          throw new Error("Token A field amount is not defined");
        if (!toFieldAmount.value)
          throw new Error("Token B field amount is not defined");

        await actions.clp.addLiquidity(
          toFieldAmount.value,
          fromFieldAmount.value
        );

        // TODO Tidy up transaction
        alert("Liquidity added");

        clearAmounts();
      },
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
        fromAmount.value = accountBalance.subtract("1").toFixed(1);
      },
      shareOfPoolPercent,
      connectedText,
    };
  },
});
</script>

<template>
  <Layout class="pool" backLink="/pool" title="Add Liquidity">
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
          toSymbolFixed
          canSwapIcon="plus"
      /></template>
      <template v-slot:default="{ requestClose }">
        <SelectTokenDialog
          :selectedTokens="[fromSymbol, toSymbol].filter(Boolean)"
          @tokenselected="requestClose"
        />
      </template>
    </Modal>

    <PriceCalculation>
      <div class="pool-share">
        <h4 class="pool-share-title text--left">Prices and pool share</h4>
        <div class="pool-share-details">
          <div v-html="aPerBRatioMessage"></div>
          <div v-html="bPerARatioMessage"></div>
          <div><span class="number">{{ shareOfPoolPercent }}</span><br>Share of Pool </div>
        </div>
      </div>
    </PriceCalculation>
    <ActionsPanel
      @nextstepclick="handleNextStepClicked"
      :nextStepAllowed="nextStepAllowed"
      :nextStepMessage="nextStepMessage"
    />
    <ModalView
      :requestClose="requestTransactionModalClose"
      :isOpen="true || transactionModalOpen"
      ><ConfirmationDialog
        @confirmswap="handleAskConfirmClicked"
        :state="transactionState"
        :requestClose="requestTransactionModalClose"
        :priceMessage="priceMessage"
        :fromToken="fromSymbol"
        :fromAmount="fromAmount"
        :toAmount="toAmount"
        :toToken="toSymbol"
    /></ModalView>
  </Layout>
</template>

<style lang="scss">
  .pool-share {
    font-size: 12px;
    font-weight: 400;
    display: flex;
    flex-direction: column;
    height: 100%;

    &-title {
      text-align: left;
      padding: 4px 16px;
      border-bottom: $divider;
    }

    &-details {
      display: flex;
      padding: 4px 16px;
      flex-grow: 1;
      justify-content: space-between;
      align-items: center;
    }
    .number {
      font-size: 16px;
      font-weight: bold;
    }
  }
</style>
