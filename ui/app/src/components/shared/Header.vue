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

<script lang="ts">
import { defineComponent } from "vue";
import Menu from "./Menu.vue";
import SifSwitch from "./SifSwitch.vue";

export default defineComponent({
  components: {
    Menu,
    SifSwitch,
  },

  data() {
    return {
      typefaceCursive: true,
    };
  },
});
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