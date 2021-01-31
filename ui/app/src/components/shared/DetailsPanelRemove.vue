<template>
  <div class="details">
    <div class="details-header">
        <h4 class="text--left">You Should Receive:</h4>

      <div class="details-row">
          <AssetItem :symbol="nativeAssetSymbol" /> {{ nativeAssetAmount }}
      </div>
      <div class="details-row">    
          <AssetItem :symbol="externalAssetSymbol" /> {{ externalAssetAmount }}
        </div>
      </div>
  </div>
</template>

<style lang="scss" scoped>
.details {
  border: 1px solid $c_gray_200;
  border-radius: $br_sm;
  background: $c_white;

  &-header {
    padding: 10px 15px;
    border-bottom: 1px solid $c_gray_200;
  }
  &-body {
    padding: 10px 15px;
  }

  &-row {
    display: flex;
    flex-direction: row;
    align-items: center;

    span:last-child {
      text-align: right;
      color: $c_gray_900;
      margin-left: auto;
    }

    span:first-child {
      color: $c_gray_700;
      font-weight: 400;
      text-align: left;
    }

    img {
      margin-right: 8px;
    }
  }
}
</style>
<script lang="ts">
import { defineComponent } from "vue";
import AssetItem from "@/components/shared/AssetItem.vue";

export default defineComponent({
  components: {
    AssetItem,
  },
  props: {
    nativeAssetSymbol: { type: String, default: ""},
    nativeAssetAmount: { type: String, default: ""},
    externalAssetAmount: { type: String, default: ""},
    externalAssetSymbol: { type: String, default: ""},
  },
});
</script>
