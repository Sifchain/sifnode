<template>
  <div v-if="show" v-bind:class="{ warning: warning, 'price-calculation': true }">
    <div class="pool-share">
      <h4 class="pool-share-title text--left"><slot name="header"></slot></h4>
      <div class="pool-share-details">
        <slot name="body"></slot>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: ["show", "warning"],
});
</script>
<style lang="scss" scoped>
.price-calculation {
  margin-bottom: 10px;
  border-radius: $br_sm;
  background: $c_white;
  overflow-y: auto;
}
.warning {
  border: 1px solid red;
}
.pool-share {
  font-size: 12px;
  font-weight: 400;
  display: flex;
  flex-direction: column;
  height: 100%;

  &-title {
    text-align: left;
    padding: 4px 16px;
    border-bottom: $divider;
  }

  &-details {
    display: flex;
    padding: 4px 16px;
    flex-grow: 1;
    justify-content: space-between;
    align-items: center;
    & > * {
      flex: 1;
      width: 33%;
    }
  }
}
</style>
