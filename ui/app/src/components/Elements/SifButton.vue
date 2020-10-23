<template>
  <button 
    class="btn"
    :class="{ 'btn--disabled': disabled, 'btn--block': block }"
    @click="$emit('click')"
  >
    <slot></slot>
  </button>
</template>

<script>

export default {
  props: {
    disabled: {
      type: Boolean,
      default: false
    },
    block: {
      type: Boolean,
      default: false
    },
  }
}
</script>

<style lang="scss">
.btn {
  @include resetButton;

  position: relative;
  display: inline-block;
  height: 30px;
  padding: 0 18px;
  overflow: hidden;
  font: inherit;
  font-size: $fs_md;
  line-height: $lh_btn;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: white;
  background: $c_gold;
  border-radius: $br_sm;
  transform: perspective(1px) translateZ(0);
  cursor: pointer;

  &:not(:last-of-type) {
    margin-right: 0.5em;
  }

  &::before {
    content: '';
    display: block;
    width: 100%;
    height: 100%;
    position: absolute;
    top: 0;
    left: 0;
    background: $g_gold;
    z-index: -1;
    opacity: 0;
    transition: opacity $trans_fast;
  }

  &:hover::before {
    opacity: 1;
  }

  &--disabled {
    cursor: default;
    background: $c_gray_400;
    color: $c_gray_800;

    &::before {
      background: transparent;
    }
  }

  &--block {
    display: block;
    width: 100%;
    margin: 0;
  }
}
</style>