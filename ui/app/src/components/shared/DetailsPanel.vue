<script lang="ts">
import { defineComponent } from "vue";
import { formatNumber, formatPercentage } from "./utils";
import Tooltip from "@/components/shared/Tooltip.vue";
import Icon from "@/components/shared/Icon.vue";

export default defineComponent({
  components: { Icon, Tooltip },
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

<template>
  <div v-if="priceMessage" class="details">
    <div class="details-header">
      <div class="details-row">
        <span>Price</span>
        <span data-handle="details-price-message">{{ priceMessage }}</span>
      </div>
    </div>
    <div
      v-if="
        (minimumReceived && toToken) || priceImpact || (providerFee && toToken)
      "
      class="details-body"
    >
      <div v-if="minimumReceived && toToken" class="details-row">
        <span>
          Minimum Received
          <Tooltip
            message="This is the minimum amount of the to token you will receive, taking into consideration the acceptable slippage percentage you are willing to take on. This amount also already takes into consideration liquidity provider fees as well. "
          >
            <Icon icon="info-box-black" />
          </Tooltip>
        </span>

        <span data-handle="details-minimum-received"
          >{{ formatNumber(minimumReceived) }}
          <span>{{
            toToken.toString().toLowerCase().includes("rowan")
              ? toToken.toString().toUpperCase()
              : "c" + toToken.slice(1).toUpperCase()
          }}</span>
        </span>
      </div>
      <div v-if="priceImpact" class="details-row">
        <span>
          Price Impact
          <Tooltip
            message="This is the percentage impact to the amount of the 'to' token in the liquidity pool based upon how much you are swapping for."
          >
            <Icon icon="info-box-black" />
          </Tooltip>
        </span>
        <span data-handle="details-price-impact">{{
          formatPercentage(priceImpact)
        }}</span>
      </div>
      <div v-if="providerFee && toToken" class="details-row">
        <span>
          Liquidity Provider Fee
          <Tooltip
            message="This is the fee paid to the liquidity providers of this pool."
          >
            <Icon icon="info-box-black" />
          </Tooltip>
        </span>
        <span data-handle="details-liquidity-provider-fee"
          >{{ showProviderFee(providerFee) }}
          <span>{{
            toToken.toString().toLowerCase().includes("rowan")
              ? toToken.toString().toUpperCase()
              : "c" + toToken.slice(1).toUpperCase()
          }}</span>
        </span>
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
    .tooltip {
      margin-left: 10px;
      display: inline-block;
    }
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
