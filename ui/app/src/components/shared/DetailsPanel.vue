<template>
  <div v-if="priceMessage" class="details">
    <div class="details-header">
      <div class="details-row">
        <span>Price</span>
        <span>{{ priceMessage }}</span>
      </div>
    </div>
    <div
      v-if="
        (minimumReceived && toToken) || priceImpact || (providerFee && toToken)
      "
      class="details-body"
    >
      <div v-if="minimumReceived && toToken" class="details-row">
        <span>Minimum Amount Received</span>
        <span
          >{{ formatNumber(minimumReceived) }}
          <span>{{ toToken.toUpperCase().replace("C", "c") }}</span></span
        >
      </div>
      <div v-if="priceImpact" class="details-row">
        <span>Price Impact</span>
        <span>{{ formatPercentage(priceImpact) }}</span>
      </div>
      <div v-if="providerFee && toToken" class="details-row">
        <span>Liquidity Provider Fee</span>
        <span
          >{{ showProviderFee(providerFee) }}
          <span>{{ toToken.toUpperCase().replace("C", "c") }}</span></span
        >
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.details {
  border: 1px solid $c_gray_200;
  border-radius: $br_sm;
  background: $c_white;

  &-header {
    padding: 10px 15px;
    border-bottom: 1px solid $c_gray_200;
  }
  &-body {
    padding: 10px 15px;
  }

  &-row {
    display: flex;
    justify-content: space-between;

    span:last-child {
      text-align: right;
      color: $c_gray_900;

      span {
        color: $c_gold_dark;
        font-weight: 700;
      }
    }

    span:first-child {
      color: $c_gray_700;
      font-weight: 400;
      text-align: left;
    }
  }
}
</style>
<script lang="ts">
import { defineComponent } from "vue";
import { formatNumber, formatPercentage } from "./utils";

export default defineComponent({
  props: {
    priceMessage: { type: String, default: "" },
    toToken: { type: String, default: "" },
    minimumReceived: { type: String, default: "" },
    providerFee: { type: String, default: "" },
    priceImpact: { type: String, default: "" },
  },
  setup() {
    function showProviderFee(providerFee: string) {
      const floatFee = parseFloat(providerFee);
      if (floatFee < 0.001) {
        return providerFee;
      } else if (floatFee < 10) {
        return floatFee.toFixed(4);
      } else if (floatFee < 100) {
        return floatFee.toFixed(3);
      } else if (floatFee < 1000) {
        return floatFee.toFixed(2);
      } else {
        return floatFee.toFixed(1);
      }
    }

    return {
      showProviderFee,
      formatNumber,
      formatPercentage,
    };
  },
});
</script>
