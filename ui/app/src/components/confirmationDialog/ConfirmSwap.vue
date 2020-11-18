<script>
import { defineComponent } from "vue";
import Loader from "@/components/shared/Loader.vue";
import SifButton from "@/components/shared/SifButton.vue";

export default defineComponent({
  components: { Loader, SifButton },
  props: {
    requestClose: Function,

  },

  data() {
    return {
      fromAmount: 125,
      toAmount: 1250,
      leastAmount: 1248.9998,
      fromToken: {
        symbol: 'usdt'
      },
      toToken: {
        symbol: 'rwn'
      },
      swapRate: 10,
      minimumReceived: 100,
      providerFee: 0.0002356,
      priceImpact: 0.134,
    }
  }
});
</script>
<template>
  <div class="confirm-swap">
    <h3 class="title mb-10">Confirm Swap</h3>
    <div class="info">
      <div class="info-row">
        <img v-if="fromToken.imageUrl" width="24" :src="fromToken.imageUrl" class="info-img" />
        <div class="placeholder" :style="backgroundStyle" v-else></div>
        <div class="info-amount">{{ fromAmount }}</div>
        <div class="info-token">{{ fromToken.symbol.toUpperCase() }}</div>
      </div>
      <div class="arrow">↓</div>
      <div class="info-row">
        <img v-if="toToken.imageUrl" width="24" :src="toToken.imageUrl" class="info-img" />
        <div class="placeholder" :style="backgroundStyle" v-else></div>
        <div class="info-amount">{{ toAmount }}</div>
        <div class="info-token">{{ toToken.symbol.toUpperCase() }}</div>
      </div>
    </div>
    <div class="estimate">Output is estimated. You will receive at least <strong>{{ leastAmount }} {{ toToken.symbol.toUpperCase() }}</strong> or the transaction will revert.</div>
    <div class="details">
      <div class="details-header">
        <div class="details-row">
          <span>Price</span>
          <span>{{ swapRate }} {{ toToken.symbol.toUpperCase() }} per {{ fromToken.symbol.toUpperCase() }}</span>
        </div>
      </div>
      <div class="details-body">
        <div class="details-row">
          <span>Minimum Received</span>
          <span>{{ minimumReceived }} {{ toToken.symbol.toUpperCase() }}</span>
        </div>
        <div class="details-row">
          <span>Price Impact</span>
          <span>{{ priceImpact }}%</span>
        </div>
        <div class="details-row">
          <span>Liquidity Provider Fee</span>
          <span>{{ providerFee }} {{ toToken.symbol.toUpperCase() }}</span>
        </div>
      </div>
    </div>
    <SifButton block primary class="confirm-btn">Confirm Swap</SifButton>
  </div>
</template>


<style lang="scss" scoped>
.confirm-swap {
  display: flex;
  flex-direction: column;
  padding: 15px 20px;
  min-height: 50vh;
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
.info {
  background: $c_gray_100;
  padding: 20px 40px 20px 20px;
  border-radius: $br_sm;

  &-row {
    display: flex;
    justify-content: start;
    align-items: center;
    font-weight: 400;
  }
  &-amount {
    font-size: 25px;
  }
  &-token {
    margin-left: auto;
    font-size: 20px;
  }
  &-img {
    margin-right: 16px;
  }
}
.arrow {
  margin: 20px 0;
  padding-left: 9px;
  text-align: left;
}
.placeholder {
  background: #aaa;
  box-sizing: border-box;
  border-radius: 16px;
  height: 24px;
  width: 24px;
  text-align: center;
  margin-right: 16px;
}

.estimate {
  margin: 25px 0;
  font-weight: 400;
  text-align: left;

  strong {
    font-weight: 700;
  }
}

.details {
  border: 1px solid $c_gray_200;
  border-radius: $br_sm;
  margin-bottom: 25px;

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

    span:first-child {
      color: $c_gray_700;
      font-weight: 400;
    }
  }
}
</style>

