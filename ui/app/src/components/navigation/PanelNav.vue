<template>
  <div class="panel-nav">
    <PanelNavItem
      to="/peg"
      color="pink"
      icon="circle-half"
      label="Peg"
    ></PanelNavItem>
    <PanelNavItem
      to="/swap"
      icon="circle-arrows"
      label="Swap"
      color="blue"
    ></PanelNavItem>
    <PanelNavItem
      to="/pool"
      icon="pool"
      label="Pool"
      color="green"
    ></PanelNavItem>
    <PanelNavItem to="/" label="Soon" disabled></PanelNavItem>
    <PanelNavItem to="/" label="Soon" disabled></PanelNavItem>
  </div>
</template>
<script>
import PanelNavItem from "./PanelNavItem.vue";
export default {
  components: {
    PanelNavItem,
  },
};
</script>
<style lang="scss" scoped>
.panel-nav {
  display: flex;
  justify-content: space-between;
}
</style>
