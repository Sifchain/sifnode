<template>
  <div>
    <div class="df fdr w100 item">
      <div class="df fdr detail-item-amount">
        <AssetItem
          v-if="copyMap[pkey].icon"
          :symbol="copyMap[pkey].icon"
          :label="false"
        />
        <!-- Spacer so numbers aligned -->
        <div v-if="!copyMap[pkey].icon" class="spacer"></div>
        <span>{{ format(item[pkey], copyMap[pkey].type) }}</span>
        <span v-if="pkey === 'multiplier'">x</span>
      </div>
      <span class="mr-3">{{ copyMap[pkey].title }}</span>
      <Tooltip :message="copyMap[pkey].tooltip">
        <Icon icon="info-box-black" />
      </Tooltip>
    </div>
  </div>
</template>

<script lang="ts">
import AssetItem from "@/components/shared/AssetItem.vue";
import Tooltip from "@/components/shared/Tooltip.vue";
import Icon from "@/components/shared/Icon.vue";

// NOTE - This will be removed and replaced with Amount API
function format(amount: number, type: string) {
  if (type === "string") {
    if (!amount) return "N/A";
    return amount;
  }
  // convert to number
  +amount;
  if (amount < 1) {
    return amount.toFixed(6);
  } else if (amount < 1000) {
    return amount.toFixed(4);
  } else if (amount < 100000) {
    return amount.toFixed(2);
  } else {
    return amount.toFixed(0);
  }
}

export default {
  components: { AssetItem, Tooltip, Icon },
  props: {
    item: { type: Object, default: null },
    copyMap: { type: Object, default: null },
    pkey: { type: String, default: null },
  },
  methods: {
    format,
  },
};
</script>

<style lang="scss" scoped>
.item {
  margin-top: 8px;
}
.detail-item-amount {
  width: 120px;
}
.spacer {
  width: 25px;
}
</style>
