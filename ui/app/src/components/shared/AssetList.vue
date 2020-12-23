<template>
  <div class="asset-list">
    <div class="line" v-for="item in items" :key="item.asset.symbol">
      <AssetItem class="token" :symbol="item.asset.symbol" />
      <div class="amount">{{ item.amount }}</div>
      <div v-if="!!item.amount" class="action">
        <SifButton primary>Peg</SifButton>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { PropType, defineComponent } from "vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import SifButton from "@/components/shared/SifButton.vue";
import { Asset } from "ui-core";
export default defineComponent({
  components: {
    AssetItem,
    SifButton,
  },
  props: {
    items: { type: Array as PropType<{ amount: string; asset: Asset }[]> },
  },
});
</script>

<style lang="scss" scoped>
.asset-list {
  background: white;
  padding: 10px;
  min-height: 300px;
  max-height: 425px;
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