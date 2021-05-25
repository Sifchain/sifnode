<script>
import { defineComponent, useCssModule } from "vue";
import SifButton from "@/components/shared/SifButton.vue";
import DetailsPanel from "@/components/shared/DetailsPanel.vue";
import AskConfirmationInfo from "@/components/shared/AskConfirmationInfo/Index.vue";

export default defineComponent({
  components: { DetailsPanel, AskConfirmationInfo, SifButton },
  props: {
    requestClose: Function,
    fromAmount: String,
    toAmount: String,
    leastAmount: String,
    fromToken: String,
    toToken: String,
    swapRate: String,
    minimumReceived: String,
    providerFee: String,
    priceImpact: String,
    priceMessage: String,
  },
  setup() {
    const styles = useCssModule();
    return { styles };
  },
});
</script>

<template>
  <div data-handle="confirm-swap-modal" :class="styles['confirm-swap']">
    <h3 class="title mb-10">Confirm Swap</h3>
    <AskConfirmationInfo :tokenAAmount="fromAmount" :tokenBAmount="toAmount" />
    <div :class="styles['estimate']">Output is estimated.</div>
    <DetailsPanel
      :class="styles['details']"
      :priceMessage="priceMessage"
      :fromToken="fromToken"
      :fromTokenImage="fromTokenImage"
      :toToken="toToken"
      :toTokenImage="toTokenImage"
      :swapRate="swapRate"
      :minimumReceived="minimumReceived"
      :providerFee="providerFee"
      :priceImpact="priceImpact"
    />
    <SifButton
      block
      primary
      :class="styles['confirm-btn']"
      @click="$emit('confirmswap')"
    >
      Confirm Swap
    </SifButton>
  </div>
</template>

<style lang="scss" module>
.confirm-swap {
  display: flex;
  flex-direction: column;
  padding: 30px 20px 20px 20px;
  min-height: 50vh;
}

.details {
  margin-bottom: 20px;
}

.confirm-btn {
  margin-top: auto !important;
}

.estimate {
  margin: 25px 0;
  font-weight: 400;
  text-align: left;

  strong {
    font-weight: 700;
  }
}
</style>
