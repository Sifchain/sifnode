

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
import WithWallet from "@/components/wallet/WithWallet.vue";
export default defineComponent({
  components: {
    CurrencyPairPanel,
    SelectTokenDialog,
    Layout,
    Modal,
    WithWallet,
  },

  setup() {
    const { api, store } = useCore();
    const marketPairFinder = api.MarketService.find;
    const fromSymbol = ref<string | null>(null);
    const fromAmount = ref<string>("0");
    const toSymbol = ref<string | null>(null);
    const toAmount = ref<string>("0");

    const selectedField = ref<"from" | "to" | null>(null);
    const { connected, connectedText } = useWalletButton({
      addrLen: 8,
    });

    const balances = computed(() => {
      return [...store.wallet.eth.balances, ...store.wallet.sif.balances];
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
      marketPairFinder,
    });

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
      handleNextStepClicked() {
        alert(
          `Swapping ${fromFieldAmount.value?.toFormatted()} for ${toFieldAmount.value?.toFormatted()}!`
        );
      },
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
    };
  },
});
</script>

<template>
  <Layout class="swap">
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
        />
      </template>
      <template v-slot:default="{ requestClose }">
        <SelectTokenDialog
          :selectedTokens="[fromSymbol, toSymbol].filter(Boolean)"
          @tokenselected="requestClose"
        />
      </template>
    </Modal>
    <div>{{ priceMessage }}</div>
    <div class="actions">
      <WithWallet>
        <template v-slot:disconnected="{ requestDialog }">
          <div class="wallet-status">No wallet connected 🅧</div>
          <button class="big-button" @click="requestDialog">
            Connect Wallet
          </button>
        </template>
        <template v-slot:connected="{ connectedText }"
          ><div>
            <div class="wallet-status">Connected to {{ connectedText }} ✅</div>
            <button
              class="big-button"
              :disabled="!nextStepAllowed"
              @click="handleNextStepClicked"
            >
              {{ nextStepMessage }}
            </button>
          </div></template
        >
      </WithWallet>
    </div>
  </Layout>
</template>

<style scoped>
.actions {
  padding-top: 1rem;
}
.big-button {
  width: 100%;
}
.wallet-status {
  margin-bottom: 1rem;
}
</style>
