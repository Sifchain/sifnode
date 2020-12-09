<script lang="ts">
import { defineComponent, ref } from "vue";
import Layout from "@/components/layout/Layout.vue";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialog from "@/components/tokenSelector/SelectTokenDialog.vue";
import Modal from "@/components/shared/Modal.vue";
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
    CurrencyPairPanel,
    SelectTokenDialog,
    PriceCalculation,
  },
  setup() {
    const { actions, poolFinder, store } = useCore();
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
      poolFinder,
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
      shareOfPoolPercent,
      connectedText,
    };
  },
});
</script>

<template>
  <Layout class="pool" backLink="/pool">
    <Modal @close="handleSelectClosed">
      <template v-slot:activator="{ requestOpen }">
        <CurrencyPairPanel
          v-model:fromAmount="fromAmount"
          v-model:fromSymbol="fromSymbol"
          @fromfocus="handleFromFocused"
          @fromblur="handleBlur"
          @fromsymbolclicked="handleFromSymbolClicked(requestOpen)"
          :fromSymbolSelectable="connected"
          v-model:toAmount="toAmount"
          v-model:toSymbol="toSymbol"
          @tofocus="handleToFocused"
          @toblur="handleBlur"
          toSymbolFixed
      /></template>
      <template v-slot:default="{ requestClose }">
        <SelectTokenDialog
          :selectedTokens="[fromSymbol, toSymbol].filter(Boolean)"
          @tokenselected="requestClose"
        />
      </template>
    </Modal>

    <PriceCalculation>
      <div>{{ aPerBRatioMessage }}</div>
      <div>{{ bPerARatioMessage }}</div>
      <div>{{ shareOfPoolPercent }}</div>
    </PriceCalculation>
    <ActionsPanel
      @nextstepclick="handleNextStepClicked"
      :nextStepAllowed="nextStepAllowed"
      :nextStepMessage="nextStepMessage"
    />
  </Layout>
</template>


