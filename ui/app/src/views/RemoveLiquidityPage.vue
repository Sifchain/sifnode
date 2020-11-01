<script lang="ts">
import { defineComponent, ref } from "vue";
import Layout from "@/components/layout/Layout.vue";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import WithWallet from "@/components/wallet/WithWallet.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialog from "@/components/tokenSelector/SelectTokenDialog.vue";
import Modal from "@/components/shared/Modal.vue";
import { PoolState, usePoolCalculator } from "../../../core/src";
import { useCore } from "@/hooks/useCore";
import { useWallet } from "@/hooks/useWallet";
import { computed } from "@vue/reactivity";
import SifButton from "@/components/shared/SifButton.vue";
import PriceCalculation from "@/components/shared/PriceCalculation.vue";

export default defineComponent({
  components: {
    Layout,
    Modal,
    CurrencyPairPanel,
    SelectTokenDialog,
    WithWallet,
    SifButton,
    PriceCalculation,
  },
  setup() {
    const { store, api } = useCore();
    const marketPairFinder = api.MarketService.find;
    const selectedField = ref<"from" | "to" | null>(null);

    const fromAmount = ref("0");
    const fromSymbol = ref<string | null>(null);
    const toAmount = ref("0");
    const toSymbol = ref<string | null>(null);

    const priceMessage = ref("");

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
      fromFieldAmount,
      toFieldAmount,
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
            return "Amount to remove is too high";
          case PoolState.VALID_INPUT:
            return "Remove Liquidity";
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
      handleNextStepClicked() {
        alert(
          `Create Pool ${fromFieldAmount.value?.toFormatted()} alongside ${toFieldAmount.value?.toFormatted()}!`
        );
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
      shareOfPool,
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
          v-model:toAmount="toAmount"
          v-model:toSymbol="toSymbol"
          @tofocus="handleToFocused"
          @toblur="handleBlur"
          @tosymbolclicked="handleToSymbolClicked(requestOpen)"
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
      <div>{{ shareOfPool }}</div>
    </PriceCalculation>
    <div class="actions">
      <WithWallet>
        <template v-slot:disconnected="{ requestDialog }">
          <div class="wallet-status">No wallet connected 🅧</div>
          <SifButton primary block @click="requestDialog">
            Connect Wallet
          </SifButton>
        </template>
        <template v-slot:connected="{ connectedText }"
          ><div>
            <div class="wallet-status">Connected to {{ connectedText }} ✅</div>
            <SifButton
              block
              primary
              :disabled="!nextStepAllowed"
              @click="handleNextStepClicked"
            >
              {{ nextStepMessage }}
            </SifButton>
          </div></template
        >
      </WithWallet>
    </div>
  </Layout>
</template>


