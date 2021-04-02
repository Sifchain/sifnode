<template>
  <div class="swap-progress">
    <div class="text"><span v-if="approving">Pending</span></div>
    <SwapProgressLoader :loading="approving" :success="approved"
      >1</SwapProgressLoader
    >
    <div class="line"></div>
    <SwapProgressLoader :loading="confirming" :success="confirmed"
      >2</SwapProgressLoader
    >
    <div class="text"><span v-if="confirming">Pending</span></div>
  </div>
</template>

<script>
import SwapProgressLoader from "./SwapProgressLoader.vue";
export default {
  components: { SwapProgressLoader },
  props: {
    approving: Boolean,
    approved: Boolean,
    confirming: Boolean,
    confirmed: Boolean,
  },
};
</script>

<style scoped lang="scss">
.swap-progress {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.line {
  height: 1px;
  width: 138px;
  margin: 0 12px;
  background-color: $c_gray_400;
}
.text {
  width: 54px;
  font-weight: 400;
  color: $c_gray_700;
}
</style>
