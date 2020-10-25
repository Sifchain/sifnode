

<template>
  <Panel class="swap-panel">
    <PanelNav />
    <div class="field-wrappers">
      <CurrencyField
        label="From"
        modelkey="from"
        @focus="handleFromFocused"
        @blur="handleBlur"
        v-model:amount="fromAmount"
        v-model:symbol="fromSymbol"
      />
      <div class="arrow">↓</div>
      <CurrencyField
        label="To"
        modelkey="to"
        @focus="handleToFocused"
        @blur="handleBlur"
        v-model:amount="toAmount"
        v-model:symbol="toSymbol"
      />
    </div>

    <div class="actions">
      <div v-if="!connected">
        <div>No wallet connected 🅧</div>
        <button class="big-button" @click="handleWalletClick">
          Connect wallet
        </button>
      </div>
      <div v-else>
        <div class="wallet-status">Connected to {{ connectedText }} ✅</div>
        <button class="big-button" :disabled="!canSwap">
          {{ nextStepMessage }}
        </button>
      </div>
    </div>
  </Panel>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { computed, ref } from "@vue/reactivity";

import { useWalletButton } from "@/components/wallet/useWalletButton";
import CurrencyField from "@/components/currencyfield/CurrencyField.vue";
import Panel from "@/components/panel/Panel.vue";
import PanelNav from "@/components/swap/PanelNav.vue";
import { useSwap } from "@/hooks/useSwap";
import { useCore } from "@/hooks/useCore";

import { useSwapCalculator } from "./swapCalculator";

export default defineComponent({
  components: { Panel, PanelNav, CurrencyField },

  setup() {
    const { api, store } = useCore();
    const marketPairFinder = api.MarketService.find;
    const {
      from: { symbol: fromSymbol, amount: fromAmount },
      to: { symbol: toSymbol, amount: toAmount },
    } = useSwap();

    const selectedField = ref<"from" | "to" | null>(null);
    const {
      connected,
      handleClicked: handleWalletClick,
      connectedText,
    } = useWalletButton({
      addrLen: 8,
    });

    const balances = computed(() => {
      return [...store.wallet.eth.balances, ...store.wallet.sif.balances];
    });

    const { nextStepMessage, priceAmount } = useSwapCalculator({
      balances,
      fromAmount,
      toAmount,
      fromSymbol,
      selectedField,
      toSymbol,
      marketPairFinder,
    });

    const canSwap = computed(() => {
      return nextStepMessage.value === "Swap"; // XXX:
    });

    function handleFromFocused() {
      selectedField.value = "from";
    }

    function handleToFocused() {
      selectedField.value = "to";
    }

    function handleBlur() {
      selectedField.value = null;
    }

    return {
      connected,
      connectedText,
      nextStepMessage,
      handleWalletClick,
      fromAmount,
      fromSymbol,
      handleFromFocused,
      handleToFocused,
      handleBlur,
      toAmount,
      toSymbol,
      canSwap,
    };
  },
});
</script>

<style scoped>
.swap-panel {
  max-width: 30rem;
}
.arrow {
  text-align: center;
  padding: 1rem;
}
.actions {
  padding-top: 1rem;
}
.big-button {
  width: 100%;
}
.wallet-status {
  margin-bottom: 1rem;
}
.field-wrappers {
  margin-bottom: 1rem;
}
</style>