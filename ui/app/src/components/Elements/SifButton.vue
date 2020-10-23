<template>
  <button 
    class="btn"
    :class="classes"
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
    medium: {
      type: Boolean,
      default: false
    },
    className: {
      type: String,
    },
    primary: {
      type: Boolean,
    },
    secondary: {
      type: Boolean,
    }
  },

  data() {
    return {
      classes: {
        'btn--disabled': this.disabled, 
        'btn--block': this.block,
        'btn--medium': this.medium,
        'btn--primary': this.primary,
        'btn--secondary': this.secondary,
        className: this.className,
      }
    }
  },
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
  border-radius: $br_sm;
  transform: perspective(1px) translateZ(0);
  cursor: pointer;

  &:not(:last-of-type) {
    margin-right: 0.5em;
  }

  &--disabled {
    cursor: default;
    pointer-events: none;
    background: $c_gray_400;
    color: $c_gray_800;

    &::before {
      background: transparent;
    }
  }

  &--primary {
    color: white;
    background: $c_gold;

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
  }

  &--secondary {
    background: $c_gray_100;
    color: $c_text;
    transition: background $trans_fast;

    &:hover {
      background: $c_gray_300;
    }
  }

  // sizes: 
  // block spans the full width of parent
  &--block {
    display: block;
    width: 100%;
    margin: 0;
  }

  &--medium {
    padding: 0 30px;
  }
}
</style>