<script>
import { defineComponent, ref } from "vue";
import { computed } from '@vue/reactivity';
import Layout from "@/components/layout/Layout.vue";
import { useAssetItem } from "@/components/shared/utils";

export default defineComponent({
  components: { Layout },
  props: { 
    pool: Object
  },
  setup(props) {
    console.log(props.pool);
    const fromSymbol = computed(() => props.pool.amounts[1].asset.symbol);
    const fromAsset = useAssetItem(fromSymbol);
    const fromToken = fromAsset.token;
    const fromBackgroundStyle = fromAsset.background;
    const fromTokenImage = computed(() => {
      if (!fromToken.value) return "";
      const t = fromToken.value;
      return t.imageUrl;
    });

    const fromValue = props.pool.amounts[1].fraction.toFixed(4);

    const toSymbol = computed(() => props.pool.amounts[0].asset.symbol);
    const toAsset = useAssetItem(toSymbol);
    const toToken = toAsset.token;
    const toBackgroundStyle = toAsset.background;
    const toTokenImage = computed(() => {
      if (!toToken.value) return "";
      const t = toToken.value;
      return t.imageUrl;
    });
    const toValue = props.pool.amounts[0].fraction.toFixed(4);

    const poolUnits = props.pool.poolUnits.toFixed(4);

    return {
      fromSymbol, fromBackgroundStyle, fromTokenImage, fromValue,
      toSymbol, toBackgroundStyle, toTokenImage, toValue,
      poolUnits,
    }
  }
});
</script>

<template>
  <Layout class="pool" @back="$emit('back')" emitBack title="Your Pair">
    <div class="sheet">
      <div class="section">
        <div class="header" @click="$emit('poolSelected')">
          <div class="image">
            <img v-if="fromTokenImage" width="22" height="22" :src="fromTokenImage" class="info-img" />
            <div class="placeholder" :style="fromBackgroundStyle" v-else></div>
            <img v-if="toTokenImage" width="22" height="22" :src="toTokenImage" class="info-img" />
            <div class="placeholder" :style="toBackgroundStyle" v-else></div>
          </div>
          <div class="symbol">
            <span>{{ fromSymbol.toUpperCase() }}</span>
            /
            <span>{{ toSymbol.toUpperCase() }}</span>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="details">

          <div class="row">
            <span>Pooled {{ fromSymbol.toUpperCase() }}:</span>
            <span class="value">{{ fromValue }}</span>
          </div>
          <div class="row">
            <span>Pooled {{ toSymbol.toUpperCase() }}:</span>
            <span class="value">{{ toValue }}</span>
          </div>

          <div class="row">
            <span>Your pool tokens:</span>
            <span>{{ poolUnits }}</span>
          </div>
        </div>
      </div>
      <div class="section">

      </div>
      <div class="section">

      </div>
    </div>
  </Layout>
</template>

<style lang="scss" scoped>
.sheet {
  background: $c_white;
  border-radius: $br_sm;
  border: $divider;

  .section {
    padding: 8px 12px;
  }

  .section:not(:last-of-type) {
    border-bottom: $divider;
  }

  .header {
    display: flex;
  }
  .symbol {
    font-size: $fs_md;
    color: $c_text;
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

  .row {
    display: flex;
    justify-content: space-between;
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
}
</style>