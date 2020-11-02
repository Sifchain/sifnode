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
import ConfirmationDialog from "@/components/confirmationDialog/ConfirmationDialog.vue";

export default defineComponent({
  components: {
    ActionsPanel,
    CurrencyPairPanel,
    Layout,
    Modal,
    PriceCalculation,
    SelectTokenDialog,
    ModalView,
    ConfirmationDialog,
  },

  setup() {
    const { api, store } = useCore();
    const marketPairFinder = api.MarketService.find;
    const fromSymbol = ref<string | null>(null);
    const fromAmount = ref<string>("0");
    const toSymbol = ref<string | null>(null);
    const toAmount = ref<string>("0");

    // const showTransactionModal = ref<boolean>(true);

    const mockTransactionState = ref<
      "selecting" | "confirming" | "confirmed" | "failed"
    >("selecting");
    const selectedField = ref<"from" | "to" | null>(null);
    const { connected, connectedText } = useWalletButton({
      addrLen: 8,
    });
    function requestTransactionModalClose() {
      mockTransactionState.value = "selecting";
    }

    const balances = computed(() => {
      return [...store.wallet.eth.balances, ...store.wallet.sif.balances];
    });

    const {
      state,
      // fromFieldAmount,
      // toFieldAmount,
      priceMessage,
    } = useSwapCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      marketPairFinder,
    });
    function handleNextStepClicked() {
      mockTransactionState.value = "confirming";
      setTimeout(() => {
        mockTransactionState.value = "confirmed";
      }, 1000);
      // alert(
      //   `Swapping ${fromFieldAmount.value?.toFormatted()} for ${toFieldAmount.value?.toFormatted()}!`
      // );
    }
    // handleNextStepClicked();
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
      nextStepAllowed: computed(() => {
        return state.value === SwapState.VALID_INPUT;
      }),
      mockTransactionState,
      transactionModalOpen: computed(() => {
        return ["confirming", "confirmed"].includes(mockTransactionState.value);
      }),
      transactionModalIsConfirmed: computed(() => {
        return mockTransactionState.value === "confirmed";
      }),
      requestTransactionModalClose,
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
            fromMax
            @fromfocus="handleFromFocused"
            @fromblur="handleBlur"
            @fromsymbolclicked="handleFromSymbolClicked(requestOpen)"
            v-model:toAmount="toAmount"
            v-model:toSymbol="toSymbol"
            @tofocus="handleToFocused"
            @toblur="handleBlur"
            @tosymbolclicked="handleToSymbolClicked(requestOpen)"
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
        ><ConfirmationDialog
          :confirmed="transactionModalIsConfirmed"
          :requestClose="requestTransactionModalClose"
      /></ModalView>
    </div>
  </Layout>
</template>