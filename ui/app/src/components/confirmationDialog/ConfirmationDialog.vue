<script lang="ts">
import { defineComponent, PropType } from "vue";

import { computed } from "@vue/reactivity";
import AskConfirmation from "./AskConfirmation.vue";
import AnimatedConfirmation from "./AnimatedConfirmation.vue";
import { ConfirmState } from "../../types";

export default defineComponent({
  components: { AskConfirmation, AnimatedConfirmation },
  inheritAttrs: false,
  props: {
    state: { type: String as PropType<ConfirmState>, default: "confirming" },
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
    transactionHash: String,
  },
  emits: ["confirmswap"],
  setup(props) {
    const confirmed = computed(() => {
      return props.state === "confirmed";
    });

    const failed = computed(() => {
      return props.state === "failed" || props.state === "rejected" || props.state === "out_of_gas";
    });

    return {
      confirmed,
      failed,
    };
  },
});
</script>
<template>
  <AskConfirmation
    v-if="state === 'confirming'"
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
    :confirmed="confirmed"
    :failed="failed"
    :state="state"
    :fromAmount="fromAmount"
    :fromToken="fromToken"
    :toAmount="toAmount"
    :toToken="toToken"
    :transactionHash="transactionHash"
    @closerequested="requestClose"
  />
</template>


