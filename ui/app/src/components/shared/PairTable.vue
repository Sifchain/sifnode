<script>
import { defineComponent } from "vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import { format } from "ui-core/src/utils/format";

export default defineComponent({
  props: ["items"],
  components: {
    AssetItem,
  },
  setup(props, context) {
    const { items } = props;
    return { items };
  },
  methods: {
    close() {
      this.opened = false;
    },
    format,
  },
});
</script>

<template>
  <div class="pairtable-container">
    <table class="pair-table" border="1">
      <tr v-for="item in items" :key="item.key">
        <td>{{ item.key }}</td>
        <td class="token">
          {{ format(item.value, { mantissa: 4 }) }}
          <AssetItem symbol="Rowan" class="ml-3" :label="false" />
        </td>
      </tr>
    </table>
  </div>
</template>

<style lang="scss" scoped>
.pair-table {
  border-spacing: 0;
  width: 100%;
  border-collapse: collapse;
  border-radius: 4px;
  border-style: hidden; /* hide standard table (collapsed) border */
  box-shadow: 0 0 0 1px #d2d2d2; /* this draws the table border  */

  td {
    border-bottom: 1px solid #d2d2d2;
    padding: 10px;
    text-align: left;
  }
  td.token {
    display: flex;
    justify-content: flex-end;
  }
  td:last-child {
    text-align: right;
  }
  tr:last-child td {
    border-bottom: none;
  }
}
</style>
