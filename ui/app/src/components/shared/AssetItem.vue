<script lang="ts">
import { computed } from "@vue/reactivity";
import { defineComponent } from "vue";
import { useAssetItem } from "./utils";

export default defineComponent({
  props: {
    symbol: String,
  },
  setup(props) {
    const symbol = computed(() => props.symbol);
    const { token, tokenLabel, backgroundStyle } = useAssetItem(symbol);

    return { token, tokenLabel, backgroundStyle };
  },
});
</script>

<template>
  <div class="row">
    <img v-if="token.imageUrl" width="16" :src="token.imageUrl" />
    <div class="placeholder" :style="backgroundStyle" v-else></div>
    <span>{{ tokenLabel }}</span>
  </div>
</template>

<style lang="scss" scoped>
.row {
  font-family: $f_default;
  display: flex;
  align-items: center;
  /* cursor: pointer; */
}

.row > * {
  margin-right: 0.5rem;
}

.placeholder {
  /* border: 3px solid #aaa; */
  background: #aaa;
  box-sizing: border-box;
  border-radius: 16px;
  height: 16px;
  width: 16px;
  text-align: center;
}
</style>