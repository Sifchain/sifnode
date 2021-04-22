<script lang="ts">
import { computed, defineComponent } from "vue";
import Icon from "@/components/shared/Icon.vue";
import Tooltip from "@/components/shared/Tooltip.vue";
import { format } from "ui-core/src/utils/format";
import { useAssetItem } from "../utils";
export default defineComponent({
  props: {
    symbol: String,
    amount: String,
  },
  components: { Icon, Tooltip },
  setup(props) {
    const symbol = computed(() => props.symbol);
    const asset = useAssetItem(symbol);

    const token = asset.token;
    const tokenLabel = asset.label.value;
    const backgroundStyle = asset.background;

    const tokenImage = computed(() => {
      if (!token.value) return "";
      const t = token.value;
      return t.imageUrl;
    });

    return { format, token, tokenLabel, tokenImage, backgroundStyle };
  },
});
</script>

<template>
  <div class="info-row">
    <img v-if="tokenImage" width="24" :src="tokenImage" class="info-img" />
    <div class="info-token">{{ tokenLabel }}</div>
    <!-- <div class="placeholder" :style="backgroundStyle" v-else></div> -->
    <div class="info-amount">{{ format(amount, { mantissa: 6 }) }}</div>
    <Tooltip :message="amount">
      <Icon icon="eye" class="info-eye" />
    </Tooltip>
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
    flex: 1 1 auto;
    text-align: right;
  }
  &-token {
    font-size: 20px;
    width: 60px;
    background: blue;
  }
  &-img {
    width: 20px;
    background: red;
    margin-right: 16px;
  }
  &-eye {
    width: 30px;
  }
}
.info-eye {
  svg {
    path {
      fill: blue !important;
    }
    fill: blue !important;
  }
}
​.placeholder {
  background: #aaa;
  box-sizing: border-box;
  border-radius: 16px;
  height: 24px;
  width: 24px;
  text-align: center;
  margin-right: 16px;
}
</style>
