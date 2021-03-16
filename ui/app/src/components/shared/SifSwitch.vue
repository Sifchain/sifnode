<template>
  <div class="sif-switch" @click="toggle">
    <div class="element" :class="{ checked: modelValue }"></div>
    <label for="" class="label">{{ label }}</label>
  </div>
</template>

<script>
export default {
  props: {
    modelValue: Boolean,
    label: String,
  },
  methods: {
    toggle: function() {
      this.$emit("update:modelValue", !this.modelValue);
    },
  },
};
</script>

<style lang="scss" scoped>
.sif-switch {
  display: flex;
  align-items: center;
  font-weight: normal;
  cursor: pointer;

  .element {
    position: relative;
    width: 50px;
    height: 22px;
    margin-right: 15px;
    border-radius: 22px;
    border: 1px solid $c_gray_700;
    background: $c_gray_700;
    transition: all $trans_fast;

    &::after {
      content: "";
      position: absolute;
      top: 0px;
      left: 0;
      display: block;
      width: 20px;
      height: 20px;
      border-radius: 20px;
      background: $c_white;
      transition: left $trans_fast;
    }

    &.checked {
      border-color: $c_gold;
      background: $c_gold;

      &::after {
        content: "";
        left: 28px;
      }
    }
  }
}
</style>
