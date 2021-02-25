<template>
  <div class="details">

    <div class="details-header">
      <div class="details-row">
        <span class="details-row-asset">
          <AssetItem :symbol="fromTokenLabel" inline />&nbsp;Deposited
        </span>
        <div class="details-row-value">
          <span>{{ fromAmount ? fromAmount : 0 }}</span>
        </div>
      </div>
      <div class="details-row">
        <span class="details-row-asset">
          <AssetItem :symbol="toTokenLabel" inline />&nbsp;Deposited
        </span>
        <div class="details-row-value">
          <span>{{ toAmount ? toAmount : 0 }}</span>
          <img
            v-if="toTokenImage"
            width="22"
            height="22"
            :src="toTokenImage"
            class="info-img"
          />
        </div>
      </div>
    </div>
    <div class="details-body">
      <div class="details-row" v-if="realBPerA">
        <span>Rates</span>
        <span>1 {{ fromTokenLabel.toLowerCase().includes("rowan") ? fromTokenLabel.toUpperCase() : "c" + fromTokenLabel.slice(1).toUpperCase() }} = {{ realBPerA }} {{ toTokenLabel.toLowerCase().includes("rowan") ? toTokenLabel.toUpperCase() : "c" + toTokenLabel.slice(1).toUpperCase() }}</span>
      </div>
      <div class="details-row" v-if="realAPerB">
        <span>&nbsp;</span>
        <span>1 {{ toTokenLabel.toLowerCase().includes("rowan") ? toTokenLabel.toUpperCase() : "c" + toTokenLabel.slice(1).toUpperCase() }} = {{ realAPerB }} {{ fromTokenLabel.toLowerCase().includes("rowan") ? fromTokenLabel.toUpperCase() : "c" + fromTokenLabel.slice(1).toUpperCase() }}</span>
      </div>
      <div class="details-row">
        <span>Share of Pool:</span>
        <span>{{ shareOfPool }}</span>
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
    justify-content: space-between;

    span:last-child {
      text-align: right;
      color: $c_gray_900;
    }

    > span:first-child {
      color: $c_gray_700;
      font-weight: 400;
      text-align: left;
    }

    &-asset {
      display: flex;
      align-items: center;
    }

    &-value {
      display: flex;
      color: $c_black;
      img {
        margin-left: 5px;
      }
    }

  }
}
</style>
<script lang="ts">
import { defineComponent } from "vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import { computed } from "@vue/reactivity";

export default defineComponent({
  components: {
    AssetItem,
  },
  props: {
    fromTokenLabel: { type: String, default: ""},
    fromAmount: { type: String, default: ""},
    fromTokenImage: { type: String, default: ""},
    toTokenLabel: { type: String, default: ""},
    toAmount: { type: String, default: ""},
    toTokenImage: { type: String, default: ""},
    aPerB: { type: String, default: ""},
    bPerA: { type: String, default: ""},
    shareOfPool: String,
  },
  setup(props) {
    const { aPerB, bPerA } = props;
    return {
      realAPerB: computed(() => {
        return aPerB === 'N/A' ? '0' : aPerB;
      }),
      realBPerA: computed(() => {
        return bPerA === 'N/A' ? '0' : bPerA;
      }),
    }
  }
});
</script>
