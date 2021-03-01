<template>
  <div class="asset-list">
    <div class="line" v-for="item in items" :key="item.asset.symbol">
      <AssetItem class="token" :symbol="item.asset.symbol" />
      <div class="amount">{{ formatNumber(item.amount) }}</div>
      <div class="action">
        <slot v-if="Number(item.amount.toFixed()) > 0" :asset="item"></slot>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { PropType, defineComponent } from "vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import { formatNumber } from "@/components/shared/utils.ts";
import { Asset } from "ui-core";
import { computed } from "@vue/reactivity";
export default defineComponent({
  components: {
    AssetItem,
  },
  props: {
    items: { type: Array as PropType<{ amount: string; asset: Asset }[]> },
  },
  methods: {
    formatNumber,
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
    text-align: right;
    margin-right: 1rem;
  }

  & .action {
    text-align: right;

    width: 100px;
  }
}
</style>
