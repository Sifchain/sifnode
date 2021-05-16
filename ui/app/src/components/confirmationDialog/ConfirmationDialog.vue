<script lang="ts">
import { defineComponent, PropType } from "vue";

import { computed, effect } from "@vue/reactivity";
import AskConfirmation from "./AskConfirmation.vue";
import AnimatedConfirmation from "./AnimatedConfirmation.vue";
// XXX: FIX THIS BEFORE PR
import { UiState } from "../../views/SwapPage.vue";
import { TransactionStatus } from "ui-core";

export default defineComponent({
  components: { AskConfirmation, AnimatedConfirmation },
  inheritAttrs: false,
  props: {
    state: {
      type: String as PropType<UiState>,
      default: "confirming",
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
    v-else
    :state="state"
    :txStatus="txStatus"
    :fromAmount="fromAmount"
    :fromToken="fromToken"
    :toAmount="toAmount"
    :toToken="toToken"
    @closerequested="requestClose"
  />
</template>
