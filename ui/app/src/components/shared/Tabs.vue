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
  <div>
    <slot></slot>
  </div>
</template>
<style lang="scss" scoped>
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
  cursor: pointer;
  user-select: none;
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
.tab-panel {
  display: none;
  &.selected {
    display: block;
  }
}
</style>
<script lang="ts">
import { defineComponent, provide, ref } from "vue";
export default defineComponent({
  emits: ["tabselected"],
  props: ["defaultIndex"],
  setup(props, context) {
    const selectedIndex = ref(props.defaultIndex || 0);
    const slots = (context.slots.default && context.slots.default()) ?? [];
    const tabs = slots.map((s: any) => ({ name: (s.props as any).title }));
    const selectedTitle = ref<string | undefined>(
      tabs[selectedIndex.value].name,
    );

    function tabSelected(index: number) {
      const selectedProps = tabs[index].name;
      selectedIndex.value = index;
      selectedTitle.value = selectedProps;
      context.emit("tabselected", {
        selectedIndex: index,
        selectedTitle: selectedProps,
      });
    }

    provide("Tabs_selectedTitle", selectedTitle);

    return {
      selectedIndex,
      tabSelected,
      tabs,
    };
  },
});
</script>
