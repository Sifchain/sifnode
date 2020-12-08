<script>
import { defineComponent } from "vue";
import SifButton from "@/components/shared/SifButton.vue";
import DetailsPanel from "@/components/shared/DetailsPanel.vue";
import AssetItemLarge, {
  getAssetLabel,
} from "@/components/shared/AssetItemLarge.vue";
import AssetItemPool from "@/components/shared/AssetItemPool.vue";
import ArrowIconButton from "@/components/shared/ArrowIconButton.vue";
import { computed } from "@vue/reactivity";

export default defineComponent({
  components: { DetailsPanel, AssetItemLarge, SifButton, ArrowIconButton, AssetItemPool },
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
  <div class="confirm-swap">
    <h3 class="title mb-10">You will receive</h3>
    <div class="info">
      
    </div>
    <div class="estimate">Output is estimated. If the price changes more than 0.5% your transaction will revert.</div>
    <div class="estimate">{{priceMessage}}</div>
    <DetailsPanel
      class="details"
      :priceMessage="priceMessage"
      :fromToken="fromToken"
      :toToken="toToken"
      :swapRate="swapRate"
      :minimumReceived="minimumReceived"
      :providerFee="providerFee"
      :priceImpact="priceImpact"
    />
    <SifButton block primary class="confirm-btn" @click="$emit('confirmswap')"
      >Confirm Supply</SifButton
    >
  </div>
</template>

<style lang="scss" scoped>
.confirm-swap {
  display: flex;
  flex-direction: column;
  padding: 30px 20px 20px 20px;
  min-height: 50vh;
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

