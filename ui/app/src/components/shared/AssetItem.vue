<script lang="ts">
import { computed, ref } from "@vue/reactivity";
import { Asset } from "ui-core";
import { defineComponent, PropType } from "vue";
import { useAssetItem } from "./utils";

export default defineComponent({
  props: {
    symbol: String,
    asset: Object as PropType<Asset>,
    inline: Boolean,
  },
  setup(props) {
    const symbol = computed(() => props.symbol);

    const asset = useAssetItem(symbol);

    const token = props.asset ? ref(props.asset) : asset.token;
    const tokenLabel = asset.label;

    const backgroundStyle = asset.background;

    return { token, tokenLabel, backgroundStyle };
  },
});
</script>

<template>
  <div class="row" :class="{'inline': inline}">
    <img v-if="token.imageUrl" width="16" :src="token.imageUrl" />
    <div class="placeholder" :style="backgroundStyle" v-else></div>
    <span>{{ tokenLabel }}</span>
  </div>
</template>

<style lang="scss" scoped>
.row {
  font-family: $f_default;
  display: flex;
  align-items: center;

  & > * {
    margin-right: 0.5rem;
  }

  &.inline {
    display: inline-flex;

    & > span {
      margin-right: 0;
    }
  }
}

.row > * {
  margin-right: 0.5rem;
}

.placeholder {
  /* border: 3px solid #aaa; */
  background: #aaa;
  box-sizing: border-box;
  border-radius: 16px;
  height: 16px;
  width: 16px;
  text-align: center;
}
</style>