<script lang="ts">
import { defineComponent, PropType, ref } from "vue";
import { computed, effect, toRefs } from "@vue/reactivity";
import Layout from "@/components/layout/Layout.vue";
import SifButton from "@/components/shared/SifButton.vue";
import { useAssetItem } from "@/components/shared/utils";
import { Fraction, LiquidityProvider, Pool, usePoolCalculator } from "ui-core";
import { useWallet } from "@/hooks/useWallet";
import { useCore } from "@/hooks/useCore";

export default defineComponent({
  components: { Layout, SifButton },
  props: {
    accountPool: Object as PropType<{
      lp: LiquidityProvider;
      pool: Pool;
    } | null>,
  },

  setup(props) {
    // TODO This needs tidying up poor componentization
    // useAssetItem should not really be used outside of a display component
    const thePool = computed(() => props.accountPool?.pool);
    const fromSymbol = computed(
      () => thePool.value?.amounts[1].asset.symbol ?? ""
    );
    const fromAsset = useAssetItem(fromSymbol);
    const fromToken = fromAsset.token;
    const fromBackgroundStyle = fromAsset.background;
    const fromTokenImage = computed(() => {
      if (!fromToken.value) return "";
      const t = fromToken.value;
      return t.imageUrl;
    });

    const fromValue = computed(() => thePool.value?.amounts[1].toFixed(0));

    const toSymbol = computed(
      () => thePool.value?.amounts[0].asset.symbol ?? ""
    );
    const toAsset = useAssetItem(toSymbol);
    const toToken = toAsset.token;
    const toBackgroundStyle = toAsset.background;
    const toTokenImage = computed(() => {
      if (!toToken.value) return "";
      const t = toToken.value;
      return t.imageUrl;
    });
    const toValue = computed(() => thePool.value?.amounts[0].toFixed(0));

    const poolUnitsAsFraction = computed(
      () => props.accountPool?.lp.units || new Fraction("0")
    );

    const myPoolShare = computed(() => {
      if (!thePool.value?.poolUnits) return null;
      const perc = poolUnitsAsFraction.value
        .divide(thePool.value?.poolUnits)
        .multiply("100")
        .toFixed(2);
      return `${perc} %`;
    });
    const myPoolUnits = computed(() => poolUnitsAsFraction.value.toFixed(0));

    return {
      fromSymbol,
      fromBackgroundStyle,
      fromTokenImage,
      fromValue,
      toSymbol,
      toBackgroundStyle,
      toTokenImage,
      toValue,
      myPoolUnits,
      myPoolShare,
    };
  },
});
</script>

<template>
  <Layout class="pool" @back="$emit('back')" emitBack title="Your Pair">
    <div class="sheet">
      <div class="section">
        <div class="header" @click="$emit('poolselected')">
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
            <span class="value">
              <span>{{ fromValue }}</span>
              <img
                v-if="fromTokenImage"
                width="22"
                height="22"
                :src="fromTokenImage"
                class="info-img"
              />
              <div
                class="placeholder"
                :style="fromBackgroundStyle"
                v-else
              ></div>
            </span>
          </div>
          <div class="row">
            <span>Pooled {{ toSymbol.toUpperCase() }}:</span>
            <span class="value">
              <span>{{ toValue }}</span>
              <img
                v-if="toTokenImage"
                width="22"
                height="22"
                :src="toTokenImage"
                class="info-img"
              />
              <div class="placeholder" :style="toBackgroundStyle" v-else></div>
            </span>
          </div>

          <div class="row">
            <span>Your pool tokens:</span>
            <span class="value">{{ myPoolUnits }}</span>
          </div>
          <div class="row">
            <span>Your pool share:</span>
            <span class="value">{{ myPoolShare }}</span>
          </div>
        </div>
      </div>

      <div class="section footer">
        <div class="mr-1">
          <div class="text--small mb-6">
            <a href="#">View pool info</a>
          </div>
          <SifButton primaryOutline nocase block>Remove Liquidity</SifButton>
        </div>
        <div class="ml-1">
          <div class="text--small mb-6">
            <a href="#">Blockexplorer</a>
          </div>
          <SifButton primary nocase block>Add Liquidity</SifButton>
        </div>
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
    padding: 2px 0;
    color: $c_text;
    font-weight: 400;

    .value {
      display: flex;
      align-items: center;
      font-weight: 700;
      & > * {
        margin-right: 0.5rem;
      }

      & > *:last-child {
        margin-right: 0;
      }
    }

    .image,
    .placeholder {
      margin-left: 4px;
    }
  }

  .info {
    text-align: left;
    font-weight: 400;
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

  .footer {
    display: flex;

    & > div {
      flex: 1;
    }
  }
}
</style>