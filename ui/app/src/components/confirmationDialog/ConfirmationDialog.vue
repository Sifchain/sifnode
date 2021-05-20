<script lang="ts">
import { defineComponent, PropType } from "vue";
import AskConfirmation from "./AskConfirmation.vue";
import AnimatedConfirmation from "./AnimatedConfirmation.vue";
import { TransactionStatus } from "ui-core";

export default defineComponent({
  components: { AskConfirmation, AnimatedConfirmation },
  inheritAttrs: false,
  props: {
    state: {
      type: String as PropType<"confirm" | "submit" | "fail" | "success">,
      default: "confirm",
    },
    txStatus: { type: Object as PropType<TransactionStatus>, default: null },
    requestClose: Function,
    priceMessage: { type: String, default: "" },
    fromAmount: String,
    fromToken: String,
    toAmount: String,
    toToken: String,
    leastAmount: String,
    swapRate: String,
    minimumReceived: String,
    providerFee: String,
    priceImpact: String,
  },
  emits: ["confirmswap"],
});
</script>
<template>
  <AskConfirmation
    v-if="state === 'confirm'"
    :fromAmount="fromAmount"
    :fromToken="fromToken"
    :toAmount="toAmount"
    :toToken="toToken"
    :leastAmount="leastAmount"
    :swapRate="swapRate"
    :minimumReceived="minimumReceived"
    :providerFee="providerFee"
    :priceImpact="priceImpact"
    :priceMessage="priceMessage"
    @confirmswap="$emit('confirmswap')"
  />
  <AnimatedConfirmation
    v-if="state === 'submit' || state === 'fail' || state === 'success'"
    :state="state"
    :txStatus="txStatus"
    :fromAmount="fromAmount"
    :fromToken="fromToken"
    :toAmount="toAmount"
    :toToken="toToken"
    @closerequested="requestClose"
  />
</template>
