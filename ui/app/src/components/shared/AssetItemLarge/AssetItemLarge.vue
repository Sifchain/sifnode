<script lang="ts">
import { computed, defineComponent } from "vue";
import Icon from "@/components/shared/Icon.vue";
import Tooltip from "@/components/shared/Tooltip.vue";
import { format } from "ui-core/src/utils/format";
import { useAssetItem } from "../utils";
import { getMantissaValue } from "ui-core/src/utils/decimalShift";

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

    return {
      format,
      token,
      tokenLabel,
      tokenImage,
      backgroundStyle,
      getMantissaValue,
    };
  },
});
</script>

<template>
  <div class="info-row" :data-handle="'info-row-' + tokenLabel.toLowerCase()">
    <img v-if="tokenImage" width="24" :src="tokenImage" class="info-img" />
    <div class="info-token">{{ tokenLabel }}</div>
    <!-- <div class="placeholder" :style="backgroundStyle" v-else></div> -->
    <div class="info-amount" data-handle="info-amount">
      {{ format(amount, { mantissa: 6 }) }}
    </div>
    <Tooltip
      v-if="getMantissaValue(amount.toString()) > 6"
      :message="amount"
      :fit="true"
    >
      <Icon icon="eye" class="info-eye" />
    </Tooltip>
    <div v-else class="eye-placeholder" />
  </div>
</template>
<style lang="scss">
.eye-placeholder {
  width: 21px;
}
.info-eye {
  svg {
    fill: #c6c6c6 !important;
    width: 16px;
    margin-left: 5px;
    margin-top: 7px;
  }
  &:hover {
    svg {
      fill: #d4b553 !important;
    }
  }
}
</style>
<style lang="scss" scoped>
.info {
  &-row {
    display: flex;
    justify-content: start;
    align-items: center;
    font-weight: 400;
  }
  &-amount {
    font-size: 18px;
    flex: 1 1 auto;
    text-align: right;
  }
  &-token {
    font-size: 18px;
    width: 60px;
  }
  &-img {
    width: 20px;
    margin-right: 10px;
  }
  &-eye {
    width: 30px;
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
