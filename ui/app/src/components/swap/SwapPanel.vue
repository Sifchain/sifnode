

<template>
  <Panel class="swap-panel">
    <PanelNav />
    <div class="field-wrappers">
      <CurrencyField label="From" v-model="fromBalance" />
      <div class="arrow">↓</div>
      <CurrencyField label="To" v-model="toBalance" />
    </div>
    <div class="info-area"></div>
    <div class="actions">
      <div v-if="!connected">
        <div>No wallet connected 🅧</div>
        <button class="big-button" @click="handleWalletClick">
          Connect wallet
        </button>
      </div>
      <div v-else>
        <div class="wallet-status">Connected to {{ connectedText }} ✅</div>
        <button class="big-button" disabled="true">Select token</button>
      </div>
    </div>
  </Panel>
</template>

<script>
import { defineComponent } from "vue";
import { useSwap } from "@/hooks/useSwap";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import CurrencyField from "@/components/currencyfield/CurrencyField.vue";
import Panel from "@/components/panel/Panel";
import PanelNav from "@/components/swap/PanelNav.vue";

export default defineComponent({
  components: { Panel, PanelNav, CurrencyField },

  setup() {
    const { swapState } = useSwap();

    const {
      connected,
      handleClicked: handleWalletClick,
      connectedText,
    } = useWalletButton({
      addrLen: 8,
    });

    return {
      connected,
      connectedText,
      fromBalance: swapState.from,
      handleWalletClick,
      toBalance: swapState.to,
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
.info-area {
  border: 1px solid#eee;
  min-height: 7rem;
}
.field-wrappers {
  margin-bottom: 1rem;
}
</style>