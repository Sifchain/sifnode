<script lang="ts">
import { defineComponent } from "vue";
import SifButton from "@/components/shared/SifButton.vue";
import DetailsPanelPool from "@/components/shared/DetailsPanelPool.vue";
import ArrowIconButton from "@/components/shared/ArrowIconButton.vue";
import { computed } from "@vue/reactivity";
import { useAssetItem } from "@/components/shared/utils";

export default defineComponent({
  components: {
    DetailsPanelPool,
    SifButton,
    // ArrowIconButton
  },
  props: {
    requestClose: Function,
    fromAmount: String,
    toAmount: String,
    leastAmount: String,
    fromToken: String,
    toToken: String,
    poolUnits: String,
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
      if (!fromAsset.token.value) return "";
      const t = fromAsset.token.value;
      return t.imageUrl;
    });

    const toSymbol = computed(() => props.toToken);
    const toAsset = useAssetItem(toSymbol);

    const toToken = toAsset.token;
    const toTokenLabel = toAsset.label;
    const toBackgroundStyle = toAsset.background;
    const toTokenImage = computed(() => {
      if (!toAsset.token.value) return "";
      const t = toAsset.token.value;
      return t.imageUrl;
    });

    return {
      fromAsset,
      // fromToken,
      fromTokenLabel,
      fromBackgroundStyle,
      fromTokenImage,
      toAsset,
      // toToken,
      toTokenLabel,
      toBackgroundStyle,
      toTokenImage,
    };
  },
});
</script>

<template>
  <div class="confirm-swap">
    <h3 class="title mb-10">Your deposit details</h3>
    <DetailsPanelPool
      class="details"
      :fromTokenLabel="fromTokenLabel"
      :toTokenLabel="toTokenLabel"
      :fromAmount="fromAmount"
      :fromTokenImage="fromTokenImage"
      :toAmount="toAmount"
      :toTokenImage="toTokenImage"
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
  margin-top: 40px;
}

.title {
  font-size: $fs_lg;
  color: $c_text;
  margin-bottom: 0;
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
