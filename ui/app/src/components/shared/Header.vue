<template>
  <header class="header">
    <div class="left">
      <Menu />
      <SifSwitch label="Typeface Cursive" v-model="typefaceCursive" />
    </div>
    <div class="right">
      <slot name="right"></slot>
    </div>
  </header>
</template>

<script>
import Menu from './Menu.vue';
import SifButton from './SifButton.vue';
import SifSwitch from './SifSwitch.vue';
import Dropdown from './Dropdown.vue';

export default {
  components: {
    Menu,
    SifButton,
    SifSwitch,
    Dropdown,
  },

  data() {
    return {
      typefaceCursive: true,
    }
  },
}
</script>

<style lang="scss" scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: $header_height;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  border-bottom: $divider;

  .left,
  .right {
    display: flex;
  }
}
</style>