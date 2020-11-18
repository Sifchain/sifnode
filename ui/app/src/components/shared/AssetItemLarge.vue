<script lang="ts">
import { computed, defineComponent } from "vue";
import { Asset, Network } from "ui-core";
import ColorHash from "color-hash";

export function getAssetLabel(t: Asset) {
  if (t.network === Network.SIFCHAIN && t.symbol.indexOf("c") === 0) {
    return ["c", t.symbol.slice(1).toUpperCase()].join("");
  }
  return t.symbol.toUpperCase();
}

export default defineComponent({
  props: {
    symbol: String,
    amount: String,
  },
  setup(props) {
    const token = computed(() =>
      props.symbol ? Asset.get(props.symbol) : undefined
    );
    const tokenLabel = computed(() => {
      if (!token.value) return "";

      return getAssetLabel(token.value);
    });
    const tokenImage = computed(() => {
      if (!token.value) return "";
      const t = token.value;
      return t.imageUrl;
    });
    const backgroundStyle = computed(() => {
      const colorHash = new ColorHash();

      const color = props.symbol ? colorHash.hex(props.symbol) : [];

      return `background: ${color};`;
    });

    return { token, tokenLabel, tokenImage, backgroundStyle };
  },
});
</script>

<template>
  <div class="info-row">
    <img v-if="tokenImage" width="24" :src="tokenImage" class="info-img" />
    <div class="placeholder" :style="backgroundStyle" v-else></div>
    <div class="info-amount">{{ amount }}</div>
    <div class="info-token">{{ tokenLabel }}</div>
  </div>
</template>

<style lang="scss" scoped>
.info {
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
.placeholder {
  background: #aaa;
  box-sizing: border-box;
  border-radius: 16px;
  height: 24px;
  width: 24px;
  text-align: center;
  margin-right: 16px;
}
</style>