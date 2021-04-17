<template>
  <div v-if="isMismatchedNetwork" class="guard">
    <h1>Easy there tiger!</h1>
    <p>
      Looks like you have a mismatched network. This means you might
      accidentally loose funds.
    </p>
    <p>
      You probably want to ensure you are connected to the Ethereum Mainnet and
      the SifChain Mainnet for any pegging and unpegging operations to work.
    </p>
  </div>
</template>

<script lang="ts">
import { computed } from "@vue/reactivity";
import { defineComponent } from "@vue/runtime-core";
import { useCore } from "../../hooks/useCore";

export default defineComponent({
  setup() {
    const { actions, store } = useCore();
    const isMismatchedNetwork = computed(() => {
      if (!store.wallet.eth.chainId || !store.wallet.sif.chainId) return false;

      return !actions.peg.isSupportedNetworkCombination(
        store.wallet.eth.chainId,
        store.wallet.sif.chainId,
      );
    });
    return { isMismatchedNetwork };
  },
});
</script>

<style scoped lang="scss">
.guard {
  position: fixed;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  background: pink;
}
</style>
