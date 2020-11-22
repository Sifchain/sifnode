<script lang="ts">
import { defineComponent, PropType } from "vue";

import { computed } from "@vue/reactivity";
import AskConfirmation from "./AskConfirmation.vue";
import AnimatedConfirmation from "./AnimatedConfirmation.vue";

export type ConfirmState =
  | "selecting"
  | "confirming"
  | "signing"
  | "confirmed"
  | "failed";

export default defineComponent({
  components: { AskConfirmation, AnimatedConfirmation },
  props: {
    state: { type: String as PropType<ConfirmState>, default: "confirming" },
    requestClose: Function,
    priceMessage: { type: String, default: "" },
  },
  setup(props) {
    const confirmed = computed(() => {
      return props.state === "confirmed";
    });

    return {
      confirmed,
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
    :state="state"
    :fromAmount="fromAmount"
    :fromToken="fromToken"
    :toAmount="toAmount"
    :toToken="toToken"
    @closerequested="requestClose"
  />
</template>


