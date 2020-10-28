

<template>
  <div class="field-wrappers">
    <CurrencyField
      label="From"
      modelkey="from"
      @focus="handleFromFocused"
      @blur="handleFromBlur"
      :amount="from.amount"
      @update:amount="handleFromUpdateAmount"
      :symbol="from.symbol"
      @update:symbol="handleFromUpdateSymbol"
    />
    <div class="arrow">↓</div>
    <CurrencyField
      label="To"
      modelkey="to"
      @focus="handleToFocused"
      @blur="handleToBlur"
      :amount="to.amount"
      @update:amount="handleToUpdateAmount"
      :symbol="to.symbol"
      @update:symbol="handleToUpdateSymbol"
    />
  </div>
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
import { defineComponent } from "vue";
import CurrencyField from "@/components/currencyfield/CurrencyField.vue";

export default defineComponent({
  components: { CurrencyField },
  props: {
    fromAmount: String,
    toAmount: String,
    fromSymbol: String,
    toSymbol: String,
  },
  emits: [
    "from-focus",
    "from-blur",
    "to-focus",
    "to-blur",
    "update:toAmount",
    "update:toSymbol",
    "update:fromAmount",
    "update:fromSymbol",
  ],
  setup(props, context) {
    function handleFromUpdateAmount(value: string) {
      context.emit("update:fromAmount", value);
    }
    function handleFromUpdateSymbol(value: string) {
      context.emit("update:fromSymbol", value);
    }

    function handleToUpdateAmount(value: string) {
      context.emit("update:toAmount", value);
    }

    function handleToUpdateSymbol(value: string) {
      context.emit("update:toSymbol", value);
    }
    function handleFromFocused() {
      context.emit("from-focus");
    }
    function handleFromBlur() {
      context.emit("from-blur");
    }
    function handleToFocused() {
      context.emit("to-focus");
    }
    function handleToBlur() {
      context.emit("to-blur");
    }
    return {
      from: { amount: props.fromAmount, symbol: props.fromSymbol },
      to: { amount: props.toAmount, symbol: props.toSymbol },
      handleFromUpdateAmount,
      handleFromUpdateSymbol,
      handleToUpdateAmount,
      handleToUpdateSymbol,
      handleFromFocused,
      handleFromBlur,
      handleToFocused,
      handleToBlur,
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