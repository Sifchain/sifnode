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
  <div v-on:click="open()" class="tooltip">
    <div class="tooltip-container" v-if="opened">
      <div class="tooltip-container-inner">
        {{message}}
      </div>
    </div>
    <slot></slot></div>
</template>

<style lang="scss" scoped>
.tooltip {
  cursor: pointer;
  position: relative;
  &-container {
    position: absolute;
    bottom: 21px;
    left: 15px;
    background: #fff;
    border: 1px solid #ccc;
    padding: 1px;
    border-radius: 10px;
    line-height: 13px;
    font-size: 11px;
    z-index: 10000;
    box-shadow: 0px 3px 10px #00000033;
    width: 210px;
    border-bottom-left-radius: 0;
    &-inner {
     border-radius: 10px;
     border-bottom-left-radius: 0;
     border: 1px solid #DEDEDE;
     border-bottom-left-radius: 0;
     padding: 1rem;
   }
  }
}
</style>
