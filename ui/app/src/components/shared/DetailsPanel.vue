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
        <span>Expected Amount Received</span>
        <span>{{ minimumReceived }} {{ toToken.toUpperCase() }}</span>
      </div>
      <div v-if="priceImpact" class="details-row">
        <span>Price Impact</span>
        <span>{{ priceImpact }}%</span>
      </div>
      <div v-if="providerFee && toToken" class="details-row">
        <span>Liquidity Provider Fee</span>
        <span>{{ providerFee }} {{ toToken.toUpperCase() }}</span>
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

export default defineComponent({
  props: {
    priceMessage: { type: String, default: "" },
    toToken: { type: String, default: "" },
    minimumReceived: { type: String, default: "" },
    providerFee: { type: String, default: "" },
    priceImpact: { type: String, default: "" },
  },
});
</script>
