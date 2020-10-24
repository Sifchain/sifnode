<template>
  <div class="row">
    <img v-if="token.imageUrl" width="16" :src="token.imageUrl" />
    <div class="placeholder" v-else></div>
    <span>{{ token.symbol.toUpperCase() }}</span>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from "vue";
import { Asset } from "../../../../core";
export default defineComponent({
  props: {
    symbol: String,
  },
  setup(props) {
    const token = computed(() =>
      props.symbol ? Asset.get(props.symbol) : undefined
    );
    return { token };
  },
});
</script>

<style scoped>
.row {
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