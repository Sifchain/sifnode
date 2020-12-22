<template>
  <div class="asset-list">
    <div v-for="token in tokens" :key="token.symbol">
      <AssetItem :symbol="token.symbol" />
    </div>
  </div>
</template>

<script>
import { defineComponent } from "vue";
import AssetItem from "@/components/shared/AssetItem.vue";
export default defineComponent({
  components: {
    AssetItem,
  },
  props: {
    tokens: { type: Array },
  },
});
</script>

<style lang="scss" scoped>
.asset-list {
  background: white;
  padding: 10px;
  min-height: 300px;
}
</style>