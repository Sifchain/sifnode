<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { computed, ref } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { SwapState, useSwapCalculator } from "../../../core";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import Modal from "@/components/shared/Modal.vue";
import SelectTokenDialog from "@/components/tokenSelector/SelectTokenDialog.vue";
import PriceCalculation from "@/components/shared/PriceCalculation.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import ModalView from "@/components/shared/ModalView.vue";
import AskConfirmation from "@/components/confirmationDialog/AskConfirmation.vue";
import SwapProgress from "@/components/swapProgress/SwapProgress.vue";
import { useCurrencyFieldState } from "@/hooks/useCurrencyFieldState";

export default defineComponent({
  components: {
    ActionsPanel,
    CurrencyPairPanel,
    Layout,
    Modal,
    PriceCalculation,
    SelectTokenDialog,
    ModalView,
    AskConfirmation,
    SwapProgress,
  },

  setup() {
    const { actions, poolFinder, store } = useCore();

    const {
      fromSymbol,
      fromAmount,
      toSymbol,
      toAmount,
    } = useCurrencyFieldState();
    const transactionState = ref<
      "selecting" | "confirming" | "confirmed" | "failed"
    >("selecting");
    const selectedField = ref<"from" | "to" | null>(null);
    const { connected, connectedText } = useWalletButton({
      addrLen: 8,
    });
    function requestTransactionModalClose() {
      transactionState.value = "selecting";
    }

    const balances = computed(() => {
      return store.wallet.sif.balances;
    });

    const {
      state,
      fromFieldAmount,
      toFieldAmount,
      priceMessage,
    } = useSwapCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      poolFinder,
    });

    function clearAmounts() {
      fromAmount.value = "0.0";
      toAmount.value = "0.0";
    }

    async function handleNextStepClicked() {
      if (!fromFieldAmount.value)
        throw new Error("from field amount is not defined");
      if (!toFieldAmount.value)
        throw new Error("to field amount is not defined");

      transactionState.value = "confirming";
      await actions.clp.swap(fromFieldAmount.value, toFieldAmount.value.asset);
      transactionState.value = "confirmed";
      clearAmounts();
    }

    return {
      connected,
      connectedText,
      nextStepMessage: computed(() => {
        switch (state.value) {
          case SwapState.SELECT_TOKENS:
            return "Select Tokens";
          case SwapState.ZERO_AMOUNTS:
            return "Please enter an amount";
          case SwapState.INSUFFICIENT_FUNDS:
            return "Insufficient Funds";
          case SwapState.VALID_INPUT:
            return "Swap";
        }
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
        selectedField.value = null;
      },

      fromAmount,
      toAmount,
      fromSymbol,
      toSymbol,
      priceMessage,
      handleFromMaxClicked() {
        selectedField.value = "from";
        const accountBalance = balances.value.find(
          (balance) => balance.asset.symbol === fromSymbol.value
        );
        if (!accountBalance) return;
        fromAmount.value = accountBalance.subtract("1").toFixed(1);
      },
      nextStepAllowed: computed(() => {
        return state.value === SwapState.VALID_INPUT;
      }),
      transactionState,
      transactionModalOpen: computed(() => {
        return ["confirming", "confirmed"].includes(transactionState.value);
      }),
      transactionModalIsConfirmed: computed(() => {
        return transactionState.value === "confirmed";
      }),
      requestTransactionModalClose,
      confirmHandler() {
        // console.log("oh My");
      },
    };
  },
});
</script>

<template>
  <Layout class="swap">
    <div>
      <Modal @close="handleSelectClosed">
        <template v-slot:activator="{ requestOpen }">
          <CurrencyPairPanel
            v-model:fromAmount="fromAmount"
            v-model:fromSymbol="fromSymbol"
            :fromMax="!!fromSymbol"
            @frommaxclicked="handleFromMaxClicked"
            @fromfocus="handleFromFocused"
            @fromblur="handleBlur"
            @fromsymbolclicked="handleFromSymbolClicked(requestOpen)"
            :fromSymbolSelectable="connected"
            v-model:toAmount="toAmount"
            v-model:toSymbol="toSymbol"
            @tofocus="handleToFocused"
            @toblur="handleBlur"
            @tosymbolclicked="handleToSymbolClicked(requestOpen)"
            :toSymbolSelectable="connected"
          />
        </template>
        <template v-slot:default="{ requestClose }">
          <SelectTokenDialog
            :selectedTokens="[fromSymbol, toSymbol].filter(Boolean)"
            @tokenselected="requestClose"
          />
        </template>
      </Modal>

      <PriceCalculation>
        {{ priceMessage }}
      </PriceCalculation>
      <ActionsPanel
        @nextstepclick="handleNextStepClicked"
        :nextStepAllowed="nextStepAllowed"
        :nextStepMessage="nextStepMessage"
      />
      <ModalView
        :requestClose="requestTransactionModalClose"
        :isOpen="transactionModalOpen"
      >
        <AskConfirmation
          :fromAmount="125"
          :fromToken="'usdt'"
          :toAmount="1250"
          :toToken="'rwn'"
          :leastAmount="1248.976"
          :swapRate="10"
          :minimumReceived="100"
          :providerFee="0.0002356"
          :priceImpact="0.134"
          @confirmswap="confirmHandler"
        />
        <!-- <ConfirmationDialog
          :confirmed="transactionModalIsConfirmed"
          :requestClose="requestTransactionModalClose"
        /> -->
      </ModalView>
    </div>

    <template v-slot:after v-if="true">
      <div>
        <SwapProgress
          :approving="false"
          :approved="true"
          :confirming="true"
          :confirmed="false"
        />
      </div>
    </template>
  </Layout>
</template>