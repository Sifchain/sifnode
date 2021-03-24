<template>
  <div class="wrapper">
    <SifButton secondary @click="toggleDropdown"><slot></slot></SifButton>
    <div class="dropdown" v-show="toggle">
      <div class="inner">
        <slot name="dropdown"></slot>
      </div>
    </div>
  </div>
</template>

<script>
import SifButton from "./SifButton.vue";

export default {
  components: {
    SifButton,
  },

  data() {
    return {
      toggle: false,
    };
  },

  methods: {
    toggleDropdown() {
      console.log("TOGGLEEE");
      this.toggle = !this.toggle;
    },

    close(e) {
      if (!this.$el.contains(e.target)) {
        this.toggle = false;
      }
    },
  },

  mounted() {
    document.addEventListener("click", this.close);
  },
  beforeDestroy() {
    document.removeEventListener("click", this.close);
  },
};
</script>

<style lang="scss" scoped>
.wrapper {
  display: inline-block;
  position: relative;
}
.dropdown {
  position: absolute;
  top: 85px;
  right: 0;
  padding: 3px;
  background: white;
  border-radius: $br_md;
  border-top-right-radius: 0;
  box-shadow: $bs_dropdown;

  .inner {
    border-radius: $br_md;
    border-top-right-radius: 0;
    border: $divider;
    padding: 1rem;
  }
}
</style>
