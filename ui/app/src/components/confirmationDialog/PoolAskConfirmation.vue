<script lang="ts">
import { defineComponent } from "vue";
import SifButton from "@/components/shared/SifButton.vue";
import DetailsPanelPool from "@/components/shared/DetailsPanelPool.vue";
import ArrowIconButton from "@/components/shared/ArrowIconButton.vue";
import { computed } from "@vue/reactivity";
import { useAssetItem } from "@/components/shared/utils";

export default defineComponent({
  components: { DetailsPanelPool, SifButton, ArrowIconButton },
  props: {
    requestClose: Function,
    fromAmount: String,
    toAmount: String,
    leastAmount: String,
    fromToken: String,
    toToken: String,
    aPerB: Number,
    bPerA: Number,
    shareOfPool: Number,
  },
  setup(props) {
    const fromSymbol = computed(() => props.fromToken);
    const fromAsset = useAssetItem(fromSymbol);

    const fromToken = fromAsset.token;
    const fromTokenLabel = fromAsset.label;
    const fromBackgroundStyle = fromAsset.background;
    const fromTokenImage = computed(() => {
      if (!fromToken.value) return "";
      const t = fromToken.value;
      return t.imageUrl;
    });

    const toSymbol = computed(() => props.toToken);
    const toAsset = useAssetItem(toSymbol);

    const toToken = toAsset.token;
    const toTokenLabel = toAsset.label;
    const toBackgroundStyle = toAsset.background;
    const toTokenImage = computed(() => {
      if (!toToken.value) return "";
      const t = toToken.value;
      return t.imageUrl;
    });

    return { 
      fromAsset, fromToken, fromTokenLabel, fromBackgroundStyle, fromTokenImage, 
      toAsset, toToken, toTokenLabel, toBackgroundStyle, toTokenImage,
    };

  },
});
</script>

<template>
  <div class="confirm-swap">
    <h3 class="title mb-10">You will receive</h3>
    <div class="pool-token">
      <div class="pool-token-value">
        <!-- TODO - what's this value? Where do I read it from? -->
        0.0000273
      </div>
      <div class="pool-token-image">
        <img v-if="fromTokenImage" width="24" :src="fromTokenImage" class="info-img" />
        <div class="placeholder" :style="fromBackgroundStyle" v-else></div>
        <img v-if="toTokenImage" width="24" :src="toTokenImage" class="info-img" />
        <div class="placeholder" :style="toBackgroundStyle" v-else></div>
      </div>
    </div>
    <div class="pool-token-label">
      {{fromTokenLabel}}/{{toTokenLabel}} Pool Tokens<br>
    </div>

    <div class="estimate">Output is estimated. If the price changes more than 0.5% your transaction will revert.</div>
    <DetailsPanelPool
      class="details"
      :fromTokenLabel="fromTokenLabel"
      :toTokenLabel="toTokenLabel"
      :fromAmount="fromAmount"
      :toAmount="toAmount"
      :aPerB="aPerB"
      :bPerA="bPerA"
      :shareOfPool="shareOfPool"
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

    & > * {
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
  .placeholder {
    display: inline-block;
    background: #aaa;
    box-sizing: border-box;
    border-radius: 16px;
    height: 24px;
    width: 24px;
    text-align: center;
  }
}
</style>

