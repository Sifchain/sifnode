<template>
  <div class="tab-header-holder">
    <div class="tab-header">
      <div
        v-for="(tab, index) in tabs"
        class="tab"
        :class="{ selected: index === selectedIndex }"
        :key="tab.name"
        @click="tabSelected(index)"
      >
        {{ tab.name }}
      </div>
    </div>
  </div>
</template>
<style lang="scss">
.tab-header-holder {
  display: flex;
}
.tab-header {
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
  display: flex;
  background: #818181;
}
.tab {
  display: block;
  padding: 4px 24px;
  &:first-child {
    border-top-left-radius: 4px;
  }
  &:last-child {
    border-top-right-radius: 4px;
  }
  background: #818181;
  color: #ffffff;
  &.selected {
    color: #818181;
    background: #ffffff;
    border-top-left-radius: 4px;
    border-top-right-radius: 4px;
  }
}
</style>
<script lang="ts">
import { ref } from "vue";
export default {
  setup() {
    const selectedIndex = ref(0);
    function tabSelected(index: number) {
      selectedIndex.value = index;
    }
    return {
      selectedIndex,
      tabSelected,
      tabs: [
        { name: "Standard" },
        { name: "Foo" },
        { name: "Bar" },
        { name: "Pegged" },
      ],
    };
  },
};
</script>