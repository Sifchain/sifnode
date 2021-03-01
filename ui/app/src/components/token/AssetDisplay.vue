<template>
  <div>{{ asset.symbol }}</div>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { Asset } from "ui-core";
export default defineComponent({
  props: { asset: Object as PropType<Asset> },
});
</script>
