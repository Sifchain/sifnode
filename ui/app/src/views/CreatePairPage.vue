<script lang="ts">
import { defineComponent, ref } from "vue";
import Layout from "@/components/layout/Layout.vue";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";

import { useWalletButton } from "@/components/wallet/useWalletButton";

export default defineComponent({
  components: { Layout, CurrencyPairPanel },
  setup() {
    const selectedField = ref<"from" | "to" | null>(null);

    const fromAmount = ref("0");
    const fromSymbol = ref<string | null>(null);
    const toAmount = ref("0");
    const toSymbol = ref<string | null>(null);

    function handleFromFocused() {
      selectedField.value = "from";
    }

    function handleToFocused() {
      selectedField.value = "to";
    }

    function handleBlur() {
      /**/
    }

    const priceMessage = ref("");

    const {
      connected,
      handleClicked: handleWalletClick,
      connectedText,
    } = useWalletButton({
      addrLen: 8,
    });

    return {
      fromAmount,
      fromSymbol,
      handleFromFocused,
      handleBlur,
      toAmount,
      toSymbol,
      handleToFocused,
      priceMessage,
      connected,
      handleWalletClick,
      connectedText,
      // canClickAction,
      // handleActionClicked,
      // nextActionMessage,
    };
  },
});
</script>

<template>
  <Layout class="pool" backLink="/pool">
    <CurrencyPairPanel
      v-model:fromAmount="fromAmount"
      v-model:fromSymbol="fromSymbol"
      @from-focus="handleFromFocused"
      @from-blur="handleBlur"
      v-model:toAmount="toAmount"
      v-model:toSymbol="toSymbol"
      @to-focus="handleToFocused"
      @to-blur="handleBlur"
    />
    <div>{{ priceMessage }}</div>
    <div class="actions">
      <div v-if="!connected">
        <div class="wallet-status">No wallet connected 🅧</div>
        <button class="big-button" @click="handleWalletClick">
          Connect wallet
        </button>
      </div>
      <div v-else>
        <div class="wallet-status">Connected to {{ connectedText }} ✅</div>
        <button
          class="big-button"
          :disabled="!canClickAction"
          @click="handleActionClicked"
        >
          {{ nextActionMessage }}
        </button>
      </div>
    </div>
  </Layout>
</template>

