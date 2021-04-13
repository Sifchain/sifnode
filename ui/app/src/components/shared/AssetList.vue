<template>
  <div class="asset-list">
    <div class="line" v-for="item in items" :key="item.asset.symbol">
      <AssetItem class="token" :symbol="item.asset.symbol" />
      <div :data-handle="item.asset.symbol + '-row-amount'" class="amount">
        {{ formatAssetAmount(item.amount) }}
        <slot name="annotation" v-bind="item"></slot>
      </div>
      <div class="action">
        <slot
          v-if="item.amount.greaterThan('0')"
          :asset="item"
          data-handle="item"
        ></slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { PropType, defineComponent } from "vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import { formatAssetAmount } from "@/components/shared/utils";
import { Asset, IAssetAmount, TransactionStatus } from "ui-core";
export default defineComponent({
  components: {
    AssetItem,
  },
  props: {
    items: {
      type: Array as PropType<
        { amount: IAssetAmount; asset: Asset; pegTxs: TransactionStatus[] }[]
      >,
    },
  },

  methods: {
    formatAssetAmount,
  },
});
</script>

<style lang="scss" scoped>
.asset-list {
  background: white;
  padding: 10px;
  min-height: 300px;
  max-height: 300px;
  overflow-y: auto;
}

.line {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;

  & .amount {
    flex-grow: 1;
    margin-right: 1rem;
    display: flex;
    justify-content: flex-end;
  }

  & .action {
    text-align: right;

    width: 100px;
  }
}
</style>
