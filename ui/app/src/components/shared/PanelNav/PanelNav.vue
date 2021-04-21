<template>
  <div class="panel-nav">
    <PanelNavItem
      data-handle="peg-page-button"
      to="/peg"
      color="pink"
      icon="circle-half"
      label="PEG"
    ></PanelNavItem>
    <PanelNavItem
      data-handle="swap-page-button"
      to="/swap"
      icon="circle-arrows"
      label="SWAP"
      color="blue"
    ></PanelNavItem>
    <PanelNavItem
      data-handle="pool-page-button"
      to="/pool"
      icon="pool"
      label="POOL"
      color="green"
    ></PanelNavItem>
    <PanelNavItem to="/rewards" label="REWARDS" icon="sif" color="gold">
    </PanelNavItem>
    <PanelNavItem to="/" label="SOON" disabled></PanelNavItem>
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
