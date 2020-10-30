<script lang="ts">
import { computed, defineComponent } from "vue";
import { Asset } from "../../../../core";
import ColorHash from "color-hash";

export default defineComponent({
  props: {
    symbol: String,
  },
  setup(props) {
    const token = computed(() =>
      props.symbol ? Asset.get(props.symbol) : undefined
    );

    const backgroundStyle = computed(() => {
      const colorHash = new ColorHash();

      const color = props.symbol ? colorHash.hex(props.symbol) : [];

      return `background: ${color};`;
    });

    return { token, backgroundStyle };
  },
});
</script>

<template>
  <div class="row">
    <img v-if="token.imageUrl" width="16" :src="token.imageUrl" />
    <div class="placeholder" :style="backgroundStyle" v-else></div>
    <span>{{ token.symbol.toUpperCase() }}</span>
  </div>
</template>

<style lang="scss" scoped>
.row {
  font-family: $f_default;
  display: flex;
  align-items: center;
  cursor: pointer;
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