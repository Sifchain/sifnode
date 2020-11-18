<script>
import { defineComponent } from "vue";
import SifButton from "@/components/shared/SifButton.vue";
import DetailsPanel from "@/components/shared/DetailsPanel.vue";

export default defineComponent({
  components: { DetailsPanel, SifButton },
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
  },
});
</script>

<template>
  <div class="confirm-swap">
    <h3 class="title mb-10">Confirm Swap</h3>
    <div class="info">
      <div class="info-row">
        <img
          v-if="fromToken.imageUrl"
          width="24"
          :src="fromToken.imageUrl"
          class="info-img"
        />
        <div class="placeholder" :style="backgroundStyle" v-else></div>
        <div class="info-amount">{{ fromAmount }}</div>
        <div class="info-token">{{ fromToken.toUpperCase() }}</div>
      </div>
      <div class="arrow">↓</div>
      <div class="info-row">
        <img
          v-if="toToken.imageUrl"
          width="24"
          :src="toToken.imageUrl"
          class="info-img"
        />
        <div class="placeholder" :style="backgroundStyle" v-else></div>
        <div class="info-amount">{{ toAmount }}</div>
        <div class="info-token">{{ toToken.toUpperCase() }}</div>
      </div>
    </div>
    <div class="estimate">
      Output is estimated. You will receive at least
      <strong>{{ leastAmount }} {{ toToken.toUpperCase() }}</strong> or the
      transaction will revert.
    </div>
    <DetailsPanel
      class="details"
      :priceMessage="'10 RWN per USDT'"
      :fromToken="fromToken"
      :toToken="toToken"
      :swapRate="swapRate"
      :minimumReceived="minimumReceived"
      :providerFee="providerFee"
      :priceImpact="priceImpact"
    />
    <SifButton block primary class="confirm-btn" @click="$emit('confirmswap')"
      >Confirm Swap</SifButton
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
</style>

