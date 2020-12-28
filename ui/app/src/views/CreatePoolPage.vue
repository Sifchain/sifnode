<script lang="ts">
import { defineComponent, ref } from "vue";
import { useRouter } from "vue-router";
import Layout from "@/components/layout/Layout.vue";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialogSif from "@/components/tokenSelector/SelectTokenDialogSif.vue";
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
    CurrencyPairPanel,
    SelectTokenDialogSif,
    PriceCalculation,
    ConfirmationDialog,
    ModalView,
  },
  setup() {
    const { actions, poolFinder, store } = useCore();
    const selectedField = ref<"from" | "to" | null>(null);
    const transactionState = ref<ConfirmState>("selecting");
    const router = useRouter();

    const {
      fromSymbol,
      fromAmount,
      toSymbol,
      toAmount,
    } = useCurrencyFieldState();

    toSymbol.value = "rowan";

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
      shareOfPool,
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

    function handleNextStepClicked() {
      if (!fromFieldAmount.value)
        throw new Error("from field amount is not defined");
      if (!toFieldAmount.value)
        throw new Error("to field amount is not defined");

      transactionState.value = "confirming";
    }

    async function handleAskConfirmClicked() {
      if (!fromFieldAmount.value)
        throw new Error("Token A field amount is not defined");
      if (!toFieldAmount.value)
        throw new Error("Token B field amount is not defined");

      transactionState.value = "signing";
      await actions.clp.addLiquidity(
        toFieldAmount.value,
        fromFieldAmount.value
      );
      transactionState.value = "confirmed";

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

      requestTransactionModalClose,

      transactionState,

      transactionModalOpen: computed(() => {
        return ["confirming", "signing", "confirmed"].includes(
          transactionState.value
        );
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
        <SelectTokenDialogSif
          :selectedTokens="[fromSymbol, toSymbol].filter(Boolean)"
          @tokenselected="requestClose"
        />
      </template>
    </Modal>

    <PriceCalculation>
      <div class="pool-share">
        <h4 class="pool-share-title text--left">Prices and pool share</h4>
        <div class="pool-share-details" v-if="fromSymbol && aPerBRatioMessage">
          <div>
            <span class="number">{{ aPerBRatioMessage }}</span
            ><br />
            <span
              >{{ fromSymbol.toUpperCase() }} per
              {{ toSymbol.toUpperCase() }}</span
            >
          </div>
          <div>
            <span class="number">{{ bPerARatioMessage }}</span
            ><br />
            <span
              >{{ toSymbol.toUpperCase() }} per
              {{ fromSymbol.toUpperCase() }}</span
            >
          </div>
          <div>
            <span class="number">{{ shareOfPoolPercent }}</span
            ><br />Share of Pool
          </div>
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
      :isOpen="transactionModalOpen"
      ><ConfirmationDialog
        @confirmswap="handleAskConfirmClicked"
        :state="transactionState"
        :requestClose="requestTransactionModalClose"
        :fromToken="fromSymbol"
        :fromAmount="fromAmount"
        :toAmount="toAmount"
        :toToken="toSymbol"
        :aPerB="aPerBRatioMessage"
        :bPerA="bPerARatioMessage"
        :shareOfPool="shareOfPoolPercent"
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

    div {
      flex: 33%;
    }
  }
  .number {
    font-size: 16px;
    font-weight: bold;
  }
}
</style>
