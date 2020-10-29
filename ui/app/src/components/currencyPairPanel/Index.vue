

<template>
  <div class="field-wrappers">
    <CurrencyField
      label="From"
      tabindex="1"
      @focus="handleFromFocused"
      @blur="handleFromBlur"
      :amount="fromAmount"
      @select-symbol="$emit('from-symbol-clicked')"
      @update:amount="handleFromUpdateAmount"
      :symbol="fromSymbol"
      @update:symbol="handleFromUpdateSymbol"
    />
    <div class="arrow">↓</div>
    <CurrencyField
      label="To"
      tabindex="2"
      @focus="handleToFocused"
      @blur="handleToBlur"
      :amount="toAmount"
      @select-symbol="$emit('to-symbol-clicked')"
      @update:amount="handleToUpdateAmount"
      :symbol="toSymbol"
      @update:symbol="handleToUpdateSymbol"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import CurrencyField from "@/components/currencyfield/CurrencyField.vue";

export default defineComponent({
  components: { CurrencyField },
  props: {
    priceMessage: String,
    fromAmount: String,
    fromSymbol: String,
    toAmount: String,
    toSymbol: String,
    connected: Boolean,
    nextStepMessage: String,
    canSwap: Boolean,
    connectedText: String,
  },
  emits: [
    "from-focus",
    "from-blur",
    "to-focus",
    "to-blur",
    "to-symbol-clicked",
    "from-symbol-clicked",
    "swap-clicked",
    "connect-clicked",
    "update:toAmount",
    "update:toSymbol",
    "update:fromAmount",
    "update:fromSymbol",
  ],
  setup(props, context) {
    function handleFromUpdateAmount(amount: string) {
      context.emit("update:fromAmount", amount);
    }
    function handleFromUpdateSymbol(symbol: string) {
      context.emit("update:fromSymbol", symbol);
    }

    function handleToUpdateAmount(amount: string) {
      context.emit("update:toAmount", amount);
    }

    function handleToUpdateSymbol(symbol: string) {
      context.emit("update:toSymbol", symbol);
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
.arrow {
  text-align: center;
  padding: 1rem;
}
.field-wrappers {
  margin-bottom: 1rem;
}
</style>