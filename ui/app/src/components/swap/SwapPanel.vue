

<template>
  <Panel class="swap-panel">
    <PanelNav />
    <div class="field-wrappers">
      <CurrencyField
        label="From"
        modelkey="from"
        v-model:amount="fromAmount"
        v-model:symbol="fromSymbol"
      />
      <div class="arrow">↓</div>
      <CurrencyField
        label="To"
        modelkey="to"
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
import { computed } from "@vue/reactivity";

import { useWalletButton } from "@/components/wallet/useWalletButton";
import CurrencyField from "@/components/currencyfield/CurrencyField.vue";
import Panel from "@/components/panel/Panel.vue";
import PanelNav from "@/components/swap/PanelNav.vue";
import { useSwap } from "@/hooks/useSwap";
import { useCore } from "@/hooks/useCore";

export default defineComponent({
  components: { Panel, PanelNav, CurrencyField },

  setup() {
    const { api, store } = useCore();
    const {
      from: { symbol: fromSymbol, amount: fromAmount },
      to: { symbol: toSymbol, amount: toAmount },
    } = useSwap();

    const {
      connected,
      handleClicked: handleWalletClick,
      connectedText,
    } = useWalletButton({
      addrLen: 8,
    });

    const marketPair = computed(() => {
      if (!fromSymbol.value || !toSymbol.value) return false;
      return api.MarketService.find(fromSymbol.value, toSymbol.value);
    });

    const balances = computed(() => {
      return [...store.wallet.eth.balances, ...store.wallet.sif.balances];
    });

    const fromHasBalance = computed(() => {
      const fromBal = balances.value.find(
        ({ asset: { symbol } }) => fromSymbol.value === symbol
      );
      return !!fromBal?.greaterThan(fromAmount.value);
    });

    const nextStepMessage = computed(() => {
      if (!marketPair.value) return "Select tokens";
      if (!fromHasBalance.value) return "Insufficient funds";
      return "Swap";
    });

    const canSwap = computed(() => {
      return marketPair.value && fromHasBalance.value;
    });

    return {
      connected,
      connectedText,
      nextStepMessage,
      handleWalletClick,
      fromAmount,
      fromSymbol,

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