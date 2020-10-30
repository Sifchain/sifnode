

<template>
  <div class="field-wrappers">
    <CurrencyField
      label="From"
      tabindex="1"
      @focus="handleFromFocused"
      @blur="handleFromBlur"
      :amount="fromAmount"
      @selectsymbol="$emit('fromsymbolclicked')"
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
      @selectsymbol="$emit('tosymbolclicked')"
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
    "fromfocus",
    "fromblur",
    "tofocus",
    "toblur",
    "tosymbolclicked",
    "fromsymbolclicked",
    "swapclicked",
    "connectclicked",
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
      context.emit("fromfocus");
    }
    function handleFromBlur() {
      context.emit("fromblur");
    }
    function handleToFocused() {
      context.emit("tofocus");
    }
    function handleToBlur() {
      context.emit("toblur");
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