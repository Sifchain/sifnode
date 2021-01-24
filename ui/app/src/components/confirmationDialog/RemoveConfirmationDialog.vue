<script lang="ts">
import { defineComponent, PropType } from "vue";

import { computed } from "@vue/reactivity";
import AskConfirmation from "./RemoveAskConfirmation.vue";
import AnimatedConfirmation from "./RemoveAnimatedConfirmation.vue";

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
    externalAssetSymbol: String,
    nativeAssetSymbol: String,
    externalAssetAmount: String,
    nativeAssetAmount: String,
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
    @confirmswap="$emit('confirmswap')"
    :externalAssetSymbol="externalAssetSymbol"
    :nativeAssetSymbol="nativeAssetSymbol"
    :externalAssetAmount="externalAssetAmount"
    :nativeAssetAmount="nativeAssetAmount"
  />
  <AnimatedConfirmation
    v-else
    :confirmed="confirmed"
    :state="state"
    :externalAssetSymbol="externalAssetSymbol"
    :nativeAssetSymbol="nativeAssetSymbol"
    :externalAssetAmount="externalAssetAmount"
    :nativeAssetAmount="nativeAssetAmount"
    :transactionHash="transactionHash"
    @closerequested="requestClose"
  />
</template>


