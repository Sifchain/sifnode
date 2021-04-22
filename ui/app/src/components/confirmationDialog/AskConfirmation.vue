<script>
import { defineComponent } from "vue";
import SifButton from "@/components/shared/SifButton.vue";
import DetailsPanel from "@/components/shared/DetailsPanel.vue";
import AssetItemLarge, {
  getAssetLabel,
} from "@/components/shared/AssetItemLarge/AssetItemLarge.vue";
import ArrowIconButton from "@/components/shared/ArrowIconButton.vue";
import { computed } from "@vue/reactivity";

export default defineComponent({
  components: { DetailsPanel, AssetItemLarge, SifButton, ArrowIconButton },
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
});
</script>

<template>
  <div data-handle="confirm-swap-modal" class="confirm-swap">
    <h3 class="title mb-10">Confirm Swap</h3>
    <div class="info">
      <AssetItemLarge :amount="fromAmount" :symbol="fromToken" />
      <ArrowIconButton left :enabled="false" />
      <AssetItemLarge :amount="toAmount" :symbol="toToken" />
    </div>
    <div class="estimate">Output is estimated.</div>
    <DetailsPanel
      class="details"
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
    <SifButton block primary class="confirm-btn" @click="$emit('confirmswap')">
      Confirm Swap
    </SifButton>
  </div>
</template>

<style lang="scss" scoped>
.confirm-swap {
  display: flex;
  flex-direction: column;
  padding: 30px 20px 20px 20px;
  min-height: 50vh;
}
.info {
  background: $c_gray_100;
  padding: 20px 20px 20px 20px;
  border-radius: $br_sm;
}

.details {
  margin-bottom: 20px;
}

.title {
  font-size: $fs_lg;
  color: $c_text;
  margin-bottom: 1rem;
  text-align: left;
}
.confirm-btn {
  margin-top: auto !important;
}

.arrow {
  padding: 5px 4px;
  text-align: left;
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
