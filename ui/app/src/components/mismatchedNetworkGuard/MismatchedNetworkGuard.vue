<template>
  <ErrorBoundary :error="isMismatchedNetwork">
    <template #fallback>
      <h1 style="font-size: 30px; line-height: 30px; padding-top: 20px">
        Easy there tiger!
      </h1>
      <h1 style="font-size: 160px; margin-top: 20px; margin-bottom: 50px">
        🐅
      </h1>
      <br />
      <p>
        Looks like you have a mismatched network. This means you might
        accidentally loose funds.
      </p>
      <br />
      <p>
        You are currently connected to
        {{ sifNetworkName }} which is not compatible with {{ evmNetworkName }}.
      </p>
      <p>
        You want to ensure you are connected to the appropriate EVM network that
        matches your SifChain Network for any pegging and unpegging operations
        to work.
      </p>
      <p v-if="suggestedEVMNetwork">
        We suggest you try the {{ suggestedEVMNetwork }}
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

export default defineComponent({
  components: { ErrorBoundary },
  setup() {
    const { actions, store } = useCore();

    const evmChainNameLookup = {
      "0x1": "Ethereum Mainnet",
      "0x3": "Ethereum Ropsten Testnet",
      "0x4": "Ethereum Rinkeby Testnet",
      "0x42": "Kovan Testnet",
      "0x539": "Local Ethereum Network",
      "0x1A4": "Goerli Testnet",
    };

    const sifChainNameLookup = {
      "sifchain-local": "Sifchain LocalNet",
      "sifchain-testnet": "Sifchain TestNet",
      "sifchain-betanet": "Sifchain BetaNet",
      "sifchain-devnet": "Sifchain DevNet",
    };

    const sifNetworkName = computed(() => {
      if (!store.wallet.sif.isConnected) return null;
      return sifChainNameLookup[
        store.wallet.sif.chainId as keyof typeof sifChainNameLookup
      ];
    });

    const evmNetworkName = computed(() => {
      if (!store.wallet.eth.isConnected) return null;
      return evmChainNameLookup[
        store.wallet.eth.chainId as keyof typeof evmChainNameLookup
      ];
    });

    const suggestedEVMNetwork = computed(() => {
      const chainId = store.wallet.sif.chainId;
      return chainId
        ? evmChainNameLookup[
            actions.peg.getSuggestedEVMNetwork(
              chainId,
            ) as keyof typeof evmChainNameLookup
          ]
        : "";
    });

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
    return {
      isMismatchedNetwork,
      sifNetworkName,
      evmNetworkName,
      suggestedEVMNetwork,
    };
  },
});
</script>
