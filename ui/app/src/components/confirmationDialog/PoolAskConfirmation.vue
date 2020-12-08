<script>
import { defineComponent } from "vue";
import SifButton from "@/components/shared/SifButton.vue";
import DetailsPanelPool from "@/components/shared/DetailsPanelPool.vue";
import AssetItemLarge, {
  getAssetLabel,
} from "@/components/shared/AssetItemLarge.vue";
import AssetItemPool from "@/components/shared/AssetItemPool.vue";
import ArrowIconButton from "@/components/shared/ArrowIconButton.vue";
import { computed } from "@vue/reactivity";
import { useAssetItem } from "@/components/shared/utils";

export default defineComponent({
  components: { DetailsPanelPool, AssetItemLarge, SifButton, ArrowIconButton, AssetItemPool },
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
    <div class="pool-token">
      <div class="pool-token-value">
        0.0000273
      </div>
      <div class="pool-token-image">
        <img src="https://via.placeholder.com/22/0000FF" width="26" height="26">
        <img src="https://via.placeholder.com/22/00ff00" width="26" height="26">
      </div>
    </div>
    <div class="pool-token-label">
      BTC/RWN Pool Tokens
    </div>

    <div class="estimate">Output is estimated. If the price changes more than 0.5% your transaction will revert.</div>
    <div class="estimate">{{priceMessage}}</div>
    <DetailsPanelPool
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
  margin-bottom: 2rem;
  text-align: left;
  font-weight: 400;
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

.pool-token {
  display: flex;
  margin-bottom: 8px;

  &-value {
    font-size: 30px;
    margin-right: 32px;
  }
  &-image {
    height: 26px;

    img {
      border-radius: 16px;

      &:nth-child(2) {
        position: relative; 
        left: -8px;
      }
    }
  }
  &-label {
    text-align: left;
    font-weight: 400;
  }
}
</style>

