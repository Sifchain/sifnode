<script lang="ts">
import { defineComponent, ref } from "vue";
import Layout from "@/components/layout/Layout.vue";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import WithWallet from "@/components/wallet/WithWallet.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialog from "@/components/tokenSelector/SelectTokenDialog.vue";
import Modal from "@/components/modal/Modal.vue";
import { usePoolCalculator } from "../../../core/src";
import { useCore } from "@/hooks/useCore";
import { useWallet } from "@/hooks/useWallet";

export default defineComponent({
  components: {
    Layout,
    Modal,
    CurrencyPairPanel,
    SelectTokenDialog,
    WithWallet,
  },
  setup() {
    const { store, api } = useCore();
    const marketPairFinder = api.MarketService.find;
    const selectedField = ref<"from" | "to" | null>(null);

    const fromAmount = ref("0");
    const fromSymbol = ref<string | null>(null);
    const toAmount = ref("0");
    const toSymbol = ref<string | null>(null);

    function handleFromFocused() {
      selectedField.value = "from";
    }

    function handleToFocused() {
      selectedField.value = "to";
    }

    function handleBlur() {
      /**/
    }

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
      handleFromFocused,
      handleBlur,
      toAmount,
      toSymbol,
      handleToFocused,
      priceMessage,
      connected,
      aPerBRatioMessage,
      bPerARatioMessage,
      nextStepMessage: "banana",
      handleFromSymbolClicked(next: () => void) {
        selectedField.value = "from";
        next();
      },
      handleToSymbolClicked(next: () => void) {
        selectedField.value = "to";
        next();
      },
      handleSelectClosed(data: string) {
        if (selectedField.value === "from") {
          fromSymbol.value = data;
        }

        if (selectedField.value === "to") {
          toSymbol.value = data;
        }
        selectedField.value = null;
      },
      shareOfPool,
      connectedText,
      // canClickAction,
      // handleActionClicked,
      // nextActionMessage,
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
          @from-focus="handleFromFocused"
          @from-blur="handleBlur"
          @from-symbol-clicked="handleFromSymbolClicked(requestOpen)"
          v-model:toAmount="toAmount"
          v-model:toSymbol="toSymbol"
          @to-focus="handleToFocused"
          @to-blur="handleBlur"
          @to-symbol-clicked="handleToSymbolClicked(requestOpen)"
      /></template>
      <template v-slot:default="{ requestClose }">
        <SelectTokenDialog @token-selected="requestClose" />
      </template>
    </Modal>
    <div>{{ aPerBRatioMessage }}</div>
    <div>{{ bPerARatioMessage }}</div>
    <div>{{ shareOfPool }}</div>
    <div class="actions">
      <WithWallet>
        <template v-slot:disconnected="{ requestDialog }">
          <div class="wallet-status">No wallet connected 🅧</div>
          <button @click="requestDialog">Connect Wallet</button>
        </template>
        <template v-slot:connected="{ connectedText }"
          ><div>
            <div class="wallet-status">Connected to {{ connectedText }} ✅</div>
            <button
              class="big-button"
              :disabled="!canSwap"
              @click="handleSwapClicked"
            >
              {{ nextStepMessage }}
            </button>
          </div></template
        >
      </WithWallet>
    </div>
  </Layout>
</template>

