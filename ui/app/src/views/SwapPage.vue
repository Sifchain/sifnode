<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { computed, ref } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { SwapState, useSwapCalculator } from "ui-core";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import Modal from "@/components/shared/Modal.vue";
import SelectTokenDialogSif from "@/components/tokenSelector/SelectTokenDialogSif.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationDialog from "@/components/confirmationDialog/ConfirmationDialog.vue";
import { useCurrencyFieldState } from "@/hooks/useCurrencyFieldState";
import DetailsPanel from "@/components/shared/DetailsPanel.vue";
import SlippagePanel from "@/components/slippagePanel/Index.vue";
import { ConfirmState } from "../types";
import { toConfirmState } from "./utils/toConfirmState";
import { Fraction } from "../../../core/src/entities";
import { getMaxAmount } from "./utils/getMaxAmount";

export default defineComponent({
  components: {
    ActionsPanel,
    CurrencyPairPanel,
    Layout,
    Modal,
    DetailsPanel,
    SelectTokenDialogSif,
    ModalView,
    ConfirmationDialog,
    SlippagePanel,
  },

  setup() {
    const { actions, poolFinder, store } = useCore();

    const {
      fromSymbol,
      fromAmount,
      toSymbol,
      toAmount,
    } = useCurrencyFieldState();

    const slippage = ref<string>("1.0");
    const transactionState = ref<ConfirmState>("selecting");
    const transactionHash = ref<string | null>(null);
    const selectedField = ref<"from" | "to" | null>(null);
    const { connected } = useWalletButton();

    function requestTransactionModalClose() {
      transactionState.value = "selecting";
    }

    const balances = computed(() => {
      return store.wallet.sif.balances;
    });

    const getAccountBalance = () => {
      return balances.value.find(
        (balance) => balance.asset.symbol === fromSymbol.value,
      );
    };

    const isFromMaxActive = computed(() => {
      const accountBalance = getAccountBalance();
      if (!accountBalance) return false;
      return fromAmount.value === accountBalance.toFixed();
    });

    const {
      state,
      fromFieldAmount,
      toFieldAmount,
      priceMessage,
      priceImpact,
      providerFee,
      minimumReceived,
    } = useSwapCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      slippage,
      poolFinder,
    });

    function clearAmounts() {
      fromAmount.value = "0.0";
      toAmount.value = "0.0";
    }

    function handleNextStepClicked() {
      if (!fromFieldAmount.value)
        throw new Error("from field amount is not defined");
      if (!toFieldAmount.value)
        throw new Error("to field amount is not defined");

      transactionState.value = "confirming";
    }

    async function handleAskConfirmClicked() {
      if (!fromFieldAmount.value)
        throw new Error("from field amount is not defined");
      if (!toFieldAmount.value)
        throw new Error("to field amount is not defined");
      if (!minimumReceived.value)
        throw new Error("minimumReceived amount is not defined");

      transactionState.value = "signing";

      const tx = await actions.clp.swap(
        fromFieldAmount.value,
        toFieldAmount.value.asset,
        minimumReceived.value,
      );
      transactionHash.value = tx.hash;
      transactionState.value = toConfirmState(tx.state); // TODO: align states
      clearAmounts();
    }

    function swapInputs() {
      selectedField.value === "to"
        ? (selectedField.value = "from")
        : (selectedField.value = "to");
      const fromAmountValue = fromAmount.value;
      const fromSymbolValue = fromSymbol.value;
      fromAmount.value = toAmount.value;
      fromSymbol.value = toSymbol.value;
      toAmount.value = fromAmountValue;
      toSymbol.value = fromSymbolValue;
    }

    return {
      connected,
      nextStepMessage: computed(() => {
        switch (state.value) {
          case SwapState.SELECT_TOKENS:
            return "Select Tokens";
          case SwapState.ZERO_AMOUNTS:
            return "Please enter an amount";
          case SwapState.INSUFFICIENT_FUNDS:
            return "Insufficient Funds";
          case SwapState.INSUFFICIENT_LIQUIDITY:
            return "Insufficient Liquidity";
          case SwapState.INVALID_AMOUNT:
            return "Invalid Amount";
          case SwapState.VALID_INPUT:
            return "Swap";
        }
      }),
      disableInputFields: computed(() => {
        return state.value === SwapState.SELECT_TOKENS;
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
      handleFromFocused() {
        selectedField.value = "from";
      },
      handleToFocused() {
        selectedField.value = "to";
      },
      handleNextStepClicked,
      handleBlur() {
        if (isFromMaxActive) return;
        selectedField.value = null;
      },
      slippage,
      fromAmount,
      toAmount,
      fromSymbol,
      minimumReceived,
      toSymbol,
      priceMessage,
      priceImpact,
      providerFee,
      handleFromMaxClicked() {
        selectedField.value = "from";
        const accountBalance = getAccountBalance();
        if (!accountBalance) return;
        const maxAmount = getMaxAmount(fromSymbol, accountBalance);
        fromAmount.value = maxAmount;
      },
      nextStepAllowed: computed(() => {
        return state.value === SwapState.VALID_INPUT;
      }),
      transactionState,
      transactionModalOpen: computed(() => {
        return [
          "confirming",
          "signing",
          "failed",
          "rejected",
          "confirmed",
          "out_of_gas",
        ].includes(transactionState.value);
      }),
      requestTransactionModalClose,
      handleArrowClicked() {
        swapInputs();
      },
      handleConfirmClicked() {
        transactionState.value = "signing";
      },
      handleAskConfirmClicked,
      transactionHash,
      isFromMaxActive,
      selectedField,
    };
  },
});
</script>

<template>
  <Layout>
    <div>
      <Modal @close="handleSelectClosed">
        <template v-slot:activator="{ requestOpen }">
          <CurrencyPairPanel
            v-model:fromAmount="fromAmount"
            v-model:fromSymbol="fromSymbol"
            :fromMax="!!fromSymbol"
            :isFromMaxActive="isFromMaxActive"
            :fromDisabled="disableInputFields"
            :toDisabled="disableInputFields"
            @frommaxclicked="handleFromMaxClicked"
            @fromfocus="handleFromFocused"
            @fromblur="handleBlur"
            @fromsymbolclicked="handleFromSymbolClicked(requestOpen)"
            :fromSymbolSelectable="connected"
            :canSwap="nextStepAllowed"
            @arrowclicked="handleArrowClicked"
            v-model:toAmount="toAmount"
            v-model:toSymbol="toSymbol"
            @tofocus="handleToFocused"
            @toblur="handleBlur"
            @tosymbolclicked="handleToSymbolClicked(requestOpen)"
            :toSymbolSelectable="connected"
            tokenALabel="From"
            tokenBLabel="To"
          />
        </template>
        <template v-slot:default="{ requestClose }">
          <SelectTokenDialogSif
            :selectedTokens="[fromSymbol, toSymbol].filter(Boolean)"
            @tokenselected="requestClose"
            :mode="selectedField"
          />
        </template>
      </Modal>
      <SlippagePanel v-if="nextStepAllowed" v-model:slippage="slippage" />
      <DetailsPanel
        :toToken="toSymbol || ''"
        :priceMessage="priceMessage || ''"
        :minimumReceived="minimumReceived || ''"
        :providerFee="providerFee || ''"
        :priceImpact="priceImpact || ''"
      />
      <ActionsPanel
        connectType="connectToSif"
        @nextstepclick="handleNextStepClicked"
        :nextStepAllowed="nextStepAllowed"
        :nextStepMessage="nextStepMessage"
      />
      <ModalView
        :requestClose="requestTransactionModalClose"
        :isOpen="transactionModalOpen"
        ><ConfirmationDialog
          @confirmswap="handleAskConfirmClicked"
          :transactionHash="transactionHash"
          :state="transactionState"
          :requestClose="requestTransactionModalClose"
          :priceMessage="priceMessage"
          :fromToken="fromSymbol"
          :fromAmount="fromAmount"
          :toAmount="toAmount"
          :toToken="toSymbol"
          :minimumReceived="minimumReceived || ''"
          :providerFee="providerFee || ''"
          :priceImpact="priceImpact || ''"
      /></ModalView>
    </div>
  </Layout>
</template>
