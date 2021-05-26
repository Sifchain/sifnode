<template>
  <div class="panel-nav">
    <PanelNavItem
      data-handle="import-page-button"
      to="/balances"
      icon="circle-half"
      label="BALANCES"
    ></PanelNavItem>
    <PanelNavItem
      data-handle="swap-page-button"
      to="/swap"
      icon="circle-arrows"
      label="SWAP"
    ></PanelNavItem>
    <PanelNavItem
      data-handle="pool-page-button"
      to="/pool"
      icon="pool"
      label="POOL"
    ></PanelNavItem>
    <PanelNavItem to="/rewards" label="REWARDS" icon="sif"> </PanelNavItem>
    <PanelNavItem
      to="/soon"
      label="SOON"
      icon="soon"
      center
      disabled
    ></PanelNavItem>
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
