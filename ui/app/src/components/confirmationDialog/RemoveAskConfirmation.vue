<script lang="ts">
import { defineComponent } from "vue";
import SifButton from "@/components/shared/SifButton.vue";
import DetailsPanelRemove from "@/components/shared/DetailsPanelRemove.vue";
import ArrowIconButton from "@/components/shared/ArrowIconButton.vue";
import { computed } from "@vue/reactivity";
import { useAssetItem } from "@/components/shared/utils";

export default defineComponent({
  components: { DetailsPanelRemove, SifButton, ArrowIconButton },
  props: {
    requestClose: Function,
    externalAssetSymbol: String,
    nativeAssetSymbol: String,
    externalAssetAmount: String,
    nativeAssetAmount: String,
  },
//   setup(props) {
//     const fromSymbol = computed(() => props.externalAssetSymbol);
//     const fromAsset = useAssetItem(fromSymbol);

//     const fromToken = fromAsset.token;
//     const fromTokenLabel = fromAsset.label;
//     const fromBackgroundStyle = fromAsset.background;
//     const fromTokenImage = computed(() => {
//       if (!fromToken.value) return "";
//       const t = fromToken.value;
//       return t.imageUrl;
//     });

//     const toSymbol = computed(() => props.nativeAssetSymbol);
//     const toAsset = useAssetItem(toSymbol);

//     const toToken = toAsset.token;
//     const toTokenLabel = toAsset.label;
//     const toBackgroundStyle = toAsset.background;
//     const toTokenImage = computed(() => {
//       if (!toToken.value) return "";
//       const t = toToken.value;
//       return t.imageUrl;
//     });

//     return { 
//       fromAsset, fromToken, fromTokenLabel, fromBackgroundStyle, fromTokenImage, 
//       toAsset, toToken, toTokenLabel, toBackgroundStyle, toTokenImage,
//     };

//   },
});
</script>

<template>
  <div class="confirm-swap">
    <h3 class="title mb-10">You are about to remove liquidity</h3>
    <DetailsPanelRemove
      class="details"
      :nativeAssetSymbol="nativeAssetSymbol"
      :externalAssetSymbol="externalAssetSymbol"
      :externalAssetAmount="externalAssetAmount"
      :nativeAssetAmount="nativeAssetAmount"
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

</style>

