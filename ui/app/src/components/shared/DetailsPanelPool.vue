<template>
  <div class="details">
    <div class="details-header">
      <div class="details-row">
        <span>BTH Deposited</span>
        <span>1.2793521</span>
      </div>
      <div class="details-row">
        <span>RWN Deposited</span>
        <span>17887,079</span>
      </div>
    </div>
    <div class="details-body">
      <div class="details-row">
        <span>Rates</span>
        <span>1 BTC = 13981.364 RWN</span>
      </div>
      <div class="details-row">
        <span>&nbsp;</span>
        <span>1 RWN = 0.0000715 BTC</span>
      </div>
      <div class="details-row">
        <span>Share of Pool:</span>
        <span>0.0043483%</span>
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
