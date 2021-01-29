<script lang="ts">
import { defineComponent, PropType } from "vue";

import { computed } from "@vue/reactivity";
import AskConfirmation from "./PoolAskConfirmation.vue";
import AnimatedConfirmation from "./PoolAnimatedConfirmation.vue";
import { Fraction } from "ui-core";
import { ConfirmState } from "../../types";

export default defineComponent({
  components: { AskConfirmation, AnimatedConfirmation },
  props: {
    state: { type: String as PropType<ConfirmState>, default: "confirming" },
    requestClose: Function,
    poolUnits: Fraction,
    fromToken: String,
    toToken: String,
    transactionHash: String,
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
    :poolUnits="poolUnits"
    :aPerB="aPerB"
    :bPerA="bPerA"
    :shareOfPool="shareOfPool"
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
    :transactionHash="transactionHash"
    @closerequested="requestClose"
  />
</template>


