

<template>
  <Panel class="swap-panel">
    <PanelNav />
    <div class="field-wrappers">
      <CurrencyField label="From" v-model="fromBalance" />
      <div class="arrow">↓</div>
      <CurrencyField label="To" v-model="toBalance" />
    </div>

    <div class="actions">
      <div v-if="!connected">
        <div>No wallet connected 🅧</div>
        <button class="big-button" @click="handleWalletClick">
          Connect wallet
        </button>
      </div>
      <div v-else>
        <div class="wallet-status">Connected to {{ connectedText }} ✅</div>
        <button class="big-button" :disabled="!canSwap">Select token</button>
      </div>
    </div>
  </Panel>
</template>

<script>
import { defineComponent } from "vue";
import { reactive, computed } from "@vue/reactivity";

import { useWalletButton } from "@/components/wallet/useWalletButton";
import CurrencyField from "@/components/currencyfield/CurrencyField.vue";
import Panel from "@/components/panel/Panel";
import PanelNav from "@/components/swap/PanelNav.vue";
import { Asset, AssetAmount, Pair } from "../../../../core";

export default defineComponent({
  components: { Panel, PanelNav, CurrencyField },

  setup() {
    const swapState = reactive({
      from: { amount: "0", symbol: null, available: null },
      to: { amount: "0", symbol: null, available: null },
    });

    const {
      connected,
      handleClicked: handleWalletClick,
      connectedText,
    } = useWalletButton({
      addrLen: 8,
    });

    const canSwap = computed(() => {
      console.log(`${swapState.from.symbol} - ${swapState.to.symbol}`);
      if (!swapState.from.symbol) return false;
      if (!swapState.to.symbol) return false;

      // Setup a new fake pools
      const ATK = Asset.get("ATK");
      const BTK = Asset.get("BTK");
      const ETH = Asset.get("ETH");

      // Setup a bunch of pairs

      const pairs = [];

      pairs.push(Pair(AssetAmount(ATK, 150), AssetAmount(BTK, 100)));
      pairs.push(Pair(AssetAmount(ATK, 100), AssetAmount(ETH, 5)));
      pairs.push(Pair(AssetAmount(BTK, 150), AssetAmount(ETH, 5)));

      const FROM = Asset.get(swapState.from.symbol);
      const TO = Asset.get(swapState.to.symbol);

      const pair = pairs.find((p) => p.contains(FROM, TO));

      return Boolean(pair);
    });

    return {
      connected,
      connectedText,
      fromBalance: swapState.from,
      handleWalletClick,
      toBalance: swapState.to,
      canSwap,
    };
  },
});
</script>

<style scoped>
.swap-panel {
  max-width: 30rem;
}
.arrow {
  text-align: center;
  padding: 1rem;
}
.actions {
  padding-top: 1rem;
}
.big-button {
  width: 100%;
}
.wallet-status {
  margin-bottom: 1rem;
}
.field-wrappers {
  margin-bottom: 1rem;
}
</style>