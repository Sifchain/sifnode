<script lang="ts">
import { computed, ref } from "@vue/reactivity";
import { defineComponent, PropType } from "vue";
import { getAssetLabel, useAssetItem } from "@/components/shared/utils";
import { LiquidityProvider, Pool } from "ui-core";

export default defineComponent({
  props: {
    accountPool: {
      type: Object as PropType<{ lp: LiquidityProvider; pool: Pool }>,
    },
  },

  setup(props) {
    const fromSymbol = computed(() =>
      props.accountPool?.pool.amounts[1].asset
        ? getAssetLabel(props.accountPool?.pool.amounts[1].asset)
        : ""
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
      props.accountPool?.pool.amounts[0].asset
        ? getAssetLabel(props.accountPool?.pool.amounts[0].asset)
        : ""
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
  <div class="pool-list-item">
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
    <div class="button">Manage</div>
  </div>
</template>

<style scoped lang="scss">
.pool-list-item {
  padding: 14px 16px;
  display: flex;

  &:not(:last-of-type) {
    border-bottom: $divider;
  }

  &:hover {
    cursor: pointer;
    background: $c_gray_50;
  }

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