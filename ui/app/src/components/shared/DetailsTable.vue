<template>
  <div v-if="header.show" class="details">
    <div class="details-header">
      <div class="details-row">
        <span>{{ header.label }}</span>
        <span>{{ header.data }}</span>
      </div>
    </div>
    <div v-if="rows.length > 0" class="details-body">
      <span v-for="row in rows" :key="row.label">
        <div v-if="row.show" class="details-row">
          <span>{{ row.label }}</span>
          <span>{{ row.data }}</span>
        </div>
      </span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.details {
  border: 1px solid $c_gray_200;
  border-radius: $br_sm;
  background: $c_white;

  &-header {
    padding: 10px 15px;
    border-bottom: 1px solid $c_gray_200;
  }
  &-body {
    padding: 10px 15px;
  }

  &-row {
    display: flex;
    justify-content: space-between;

    span:last-child {
      text-align: right;
      color: $c_gray_900;
    }

    span:first-child {
      color: $c_gray_700;
      font-weight: 400;
      text-align: left;
    }
  }
}
</style>
<script lang="ts">
import { defineComponent, PropType } from "vue";
import { computed } from "@vue/reactivity";

type Row = { show: boolean; label: string; data: string };

export default defineComponent({
  props: {
    header: {
      type: Object as PropType<Row>,
      default: () => ({ show: false, label: "", data: "" }),
    },
    rows: {
      type: Array as PropType<Row[]>,
      default: [],
    },
  },
});
</script>
