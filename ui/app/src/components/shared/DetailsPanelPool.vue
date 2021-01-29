<template>
  <div class="details">
    <div class="details-header">
      <div class="details-row">
        <span>{{ fromTokenLabel }} Deposited</span>
        <div class="details-row-value">
          <span>{{ fromAmount }}</span>
          <img
            v-if="fromTokenImage"
            width="22"
            height="22"
            :src="fromTokenImage"
            class="info-img"
          />
        </div>
      </div>
      <div class="details-row">
        <span>ROWAN Deposited</span>
        <div class="details-row-value">
          <span>{{ toAmount }}</span>
          <img
            v-if="toTokenImage"
            width="22"
            height="22"
            :src="toTokenImage"
            class="info-img"
          />
        </div>
      </div>
    </div>
    <div class="details-body">
      <div class="details-row">
        <span>Rates</span>
        <span>1 {{ fromTokenLabel }} = {{ aPerB }} {{ toTokenLabel }}</span>
      </div>
      <div class="details-row">
        <span>&nbsp;</span>
        <span>1 {{ toTokenLabel }} = {{ bPerA }} {{ fromTokenLabel }}</span>
      </div>
      <div class="details-row">
        <span>Share of Pool:</span>
        <span>{{ shareOfPool }}</span>
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

    > span:first-child {
      color: $c_gray_700;
      font-weight: 400;
      text-align: left;
    }

    &-value {
      display: flex;
      color: $c_black;
      img {
        margin-left: 5px;
      }
    }

  }
}
</style>
<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    fromTokenLabel: { type: String, default: ""},
    fromAmount: { type: Number, default: ""},
    fromTokenImage: { type: String, default: ""},
    toTokenLabel: { type: String, default: ""},
    toAmount: { type: Number, default: ""},
    toTokenImage: { type: String, default: ""},
    aPerB: { type: Number, default: ""},
    bPerA: { type: Number, default: ""},
    shareOfPool: Number,
  },
});
</script>
