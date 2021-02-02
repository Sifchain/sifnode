<script>
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    message: {
      type: String,
    },
  },

  data: function() {
  	return {
  		opened: false
    };
  },
  methods: {
    close() {
      console.log("close tool tips");
      this.opened = false
      document.body.removeEventListener("mousedown", this.close);
    },
    open() {
      this.opened = true;
      // add click handler to whole page to close tooltip
      document.body.addEventListener("mousedown", this.close);
    }
  }
});
</script>

<template>
  <div v-on:click="open()" class="tooltip"><div class="tooltip-container" v-if="opened">{{message}}</div><slot></slot></div>
</template>

<style lang="scss" scoped>
.tooltip {
  cursor: pointer;
  position: relative;
  &-container {
    position: absolute;
    bottom: 25px;
    left: -10px;
    background: #fff;
    border: 1px solid #ccc;
    padding: 15px;
    border-radius: 6px;
    line-height: 13px;
    font-size: 11px;
    z-index: 10000;
    box-shadow: 0.1px 0.1px #999;
    width: 240px;
  }
}
</style>
