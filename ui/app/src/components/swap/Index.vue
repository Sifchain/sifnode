

<template>
  <CurrencyPairPanel
    v-model:fromAmount="fromAmount"
    v-model:fromSymbol="fromSymbol"
    @from-focus="handleFromFocused"
    @from-blur="handleBlur"
    v-model:toAmount="toAmount"
    v-model:toSymbol="toSymbol"
    @to-focus="handleToFocused"
    @to-blur="handleBlur"
  />
  <div>{{ priceMessage }}</div>
  <div class="actions">
    <div v-if="!connected">
      <div class="wallet-status">No wallet connected 🅧</div>
      <button class="big-button" @click="handleWalletClick">
        Connect wallet
      </button>
    </div>
    <div v-else>
      <div class="wallet-status">Connected to {{ connectedText }} ✅</div>
      <button
        class="big-button"
        :disabled="!canSwap"
        @click="handleSwapClicked"
      >
        {{ nextStepMessage }}
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, ref } from "@vue/reactivity";
import { defineComponent } from "vue";
import { useCore } from "@/hooks/useCore";
import { useSwap } from "@/hooks/useSwap";
import { useSwapCalculator } from "./swapCalculator";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";

export default defineComponent({
  components: { CurrencyPairPanel },

  setup() {
    const { api, store } = useCore();
    const marketPairFinder = api.MarketService.find;
    const swapState = useSwap();
    const {
      from: { symbol: fromSymbol, amount: fromAmount },
      to: { symbol: toSymbol, amount: toAmount },
    } = swapState;

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

    const {
      nextStepMessage,
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

    const canSwap = computed(() => {
      return nextStepMessage.value === "Swap";
    });

    function handleFromFocused() {
      selectedField.value = "from";
    }

    function handleToFocused() {
      selectedField.value = "to";
    }

    function handleSwapClicked() {
      alert(
        `Swapping ${fromFieldAmount.value?.toFormatted()} for ${toFieldAmount.value?.toFormatted()}!`
      );
    }

    function handleBlur() {
      selectedField.value = null;
    }

    return {
      connected,
      connectedText,
      nextStepMessage,
      handleWalletClick,
      handleFromFocused,
      handleToFocused,
      handleSwapClicked,
      handleBlur,
      fromAmount: swapState.from.amount,
      toAmount: swapState.to.amount,
      fromSymbol: swapState.from.symbol,
      toSymbol: swapState.to.symbol,
      priceMessage,
      canSwap,
    };
  },
});
</script>

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