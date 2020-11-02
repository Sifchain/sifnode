<template>
  <div class="panel-nav">
    <PanelNavItem to="/swap" icon="arrows" label="Swap"></PanelNavItem>
    <PanelNavItem to="/pool" icon="circle-half" label="Pool"></PanelNavItem>
    <PanelNavItem to="/" label="Soon" disabled></PanelNavItem>
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