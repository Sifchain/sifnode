<template>
  <ErrorBoundary :error="isMismatchedNetwork">
    <template #fallback>
      <h1>Easy there tiger!</h1>
      <p>
        Looks like you have a mismatched network. This means you might
        accidentally loose funds.
      </p>
      <p>
        You probably want to ensure you are connected to the Ethereum Mainnet
        and the SifChain Mainnet for any pegging and unpegging operations to
        work.
      </p>
    </template>
    <template #default>
      <slot></slot>
    </template>
  </ErrorBoundary>
</template>

<script lang="ts">
import { computed } from "@vue/reactivity";
import { defineComponent } from "@vue/runtime-core";
import { useCore } from "../../hooks/useCore";
import ErrorBoundary from "@/components/shared/ErrorBoundary/ErrorBoundary.vue";
import { useRouter } from "vue-router";

export default defineComponent({
  components: { ErrorBoundary },
  setup() {
    const { actions, store } = useCore();

    const isMismatchedNetwork = computed(() => {
      if (
        !store.wallet.eth.isConnected ||
        !store.wallet.sif.isConnected ||
        !store.wallet.eth.chainId ||
        !store.wallet.sif.chainId
      )
        return false;

      return !actions.peg.isSupportedNetworkCombination(
        store.wallet.eth.chainId,
        store.wallet.sif.chainId,
      );
    });

    // console.log({ isMismatchedNetwork: isMismatchedNetwork.value });
    return { isMismatchedNetwork };
  },
});
</script>
