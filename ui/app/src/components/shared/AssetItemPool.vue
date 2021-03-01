<script lang="ts">
import { computed } from "@vue/reactivity";
import { LiquidityProvider, Pool } from "ui-core";
import { getAssetLabel, useAssetItem } from "@/components/shared/utils";
import { defineComponent, PropType } from "vue";

export default defineComponent({
  props: {
    pool: {
      type: Object as PropType<{ lp: LiquidityProvider; pool: Pool }>,
    },
    tokenASymbol: {
      type: String,
      default: "",
    },
    tokenBSymbol: {
      type: String,
      default: "",
    },
    inline: Boolean,
  },

  setup(props) {
    const fromSymbol = computed(() =>
      props.pool?.pool.amounts[1].asset
        ? getAssetLabel(props.pool?.pool.amounts[1].asset)
        : props.tokenASymbol
    );

    const fromAsset = useAssetItem(fromSymbol);
    const fromToken = fromAsset.token;
    const fromBackgroundStyle = fromAsset.background;
    const fromTokenImage = computed(() => {
      if (!fromToken.value) return "";
      const t = fromToken.value;
      return t.imageUrl;
    });

    const toSymbol = computed(() =>
      props.pool?.pool.amounts[0].asset
        ? getAssetLabel(props.pool?.pool.amounts[0].asset)
        : props.tokenBSymbol
    );
    const toAsset = useAssetItem(toSymbol);
    const toToken = toAsset.token;
    const toBackgroundStyle = toAsset.background;
    const toTokenImage = computed(() => {
      if (!toToken.value) return "";
      const t = toToken.value;
      return t.imageUrl;
    });

    return {
      fromSymbol,
      fromBackgroundStyle,
      fromTokenImage,
      toSymbol,
      toBackgroundStyle,
      toTokenImage,
    };
  },
});
</script>

<template>
  <div class="pool-asset" :class="{ inline: inline }">
    <div class="image">
      <img
        v-if="fromTokenImage"
        width="22"
        height="22"
        :src="fromTokenImage"
        class="info-img"
      />
      <div class="placeholder" :style="fromBackgroundStyle" v-else></div>
      <img
        v-if="toTokenImage"
        width="22"
        height="22"
        :src="toTokenImage"
        class="info-img"
      />
      <div class="placeholder" :style="toBackgroundStyle" v-else></div>
    </div>
    <div class="symbol">
      <span>{{ fromSymbol }}</span>
      /
      <span>{{ toSymbol }}</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.pool-asset {
  display: flex;

  .image {
    height: 22px;

    & > * {
      border-radius: 16px;

      &:nth-child(2) {
        position: relative;
        left: -6px;
      }
    }
  }

  &.inline {
    display: inline-flex;

    & > span {
      margin-right: 0;
    }
  }

  .placeholder {
    display: inline-block;
    background: #aaa;
    box-sizing: border-box;
    border-radius: 16px;
    height: 22px;
    width: 22px;
    text-align: center;
  }

  .symbol {
    font-size: $fs_md;
    color: $c_text;
  }

  .button {
    font-size: $fs_sm;
    font-weight: normal;
    flex-grow: 1;
    text-align: right;
  }
}
</style>
