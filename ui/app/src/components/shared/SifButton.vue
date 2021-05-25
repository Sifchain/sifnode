<template>
  <button
    v-if="!to"
    v-bind="$attrs"
    class="btn"
    :class="classes"
    :disabled="disabled"
  >
    <span class="content">
      <slot></slot>
    </span>
  </button>

  <router-link
    v-else
    v-bind="$attrs"
    :to="to"
    class="btn"
    :class="classes"
    :disabled="disabled"
  >
    <span class="content">
      <slot></slot>
    </span>
  </router-link>
</template>

<script lang="ts">
import { PropType, defineComponent } from "vue";
export default defineComponent({
  inheritAttrs: false,
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
    block: {
      type: Boolean,
      default: false,
    },
    connect: {
      type: Boolean,
      default: false,
    },
    active: {
      type: Boolean,
      default: false,
    },
    medium: {
      type: Boolean,
      default: false,
    },
    className: {
      type: String,
    },
    primary: {
      type: Boolean,
    },
    primaryOutline: {
      type: Boolean,
    },
    secondary: {
      type: Boolean,
    },
    ghost: {
      type: Boolean,
    },
    small: {
      type: Boolean,
    },
    success: {
      type: Boolean,
    },
    round: {
      type: Boolean,
    },
    nocase: {
      type: Boolean,
      default: false,
    },
    to: {
      type: String,
      default: "",
    },
  },

  data() {
    return {
      classes: {
        block: this.block,
        connect: this.connect,
        active: this.active,
        medium: this.medium,
        primary: this.primary,
        "primary-outline": this.primaryOutline,
        secondary: this.secondary,
        className: this.className,
        ghost: this.ghost,
        small: this.small,
        nocase: this.nocase,
        success: this.success,
        round: this.round,
      },
    };
  },
});
</script>

<style lang="scss" scoped>
.btn {
  @include resetButton;
  position: relative;
  display: flex;
  height: 30px;
  padding: 0 18px;
  align-items: center;
  overflow: hidden;
  font-size: $fs;
  font: inherit;
  // line-height: $lh_btn;
  letter-spacing: 0px;
  border-radius: $br_sm;
  transform: perspective(1px) translateZ(0);
  cursor: pointer;
  box-sizing: border-box;
  justify-content: center;

  &.nocase {
    text-transform: none;
    letter-spacing: 0;
    font-weight: 400;
    font-size: $fs_sm;
  }

  &:not(:last-of-type) {
    margin-right: 0.5em;
  }

  &:disabled,
  &[disabled="true"] {
    cursor: none;
    pointer-events: none;
    background: $c_gray_400 !important;
    color: $c_gray_800 !important;

    &::before {
      background: transparent;
    }
  }

  &.primary {
    color: white;
    background: $c_gold;
    font-weight: initial;
    font-style: normal;
    &::before {
      content: "";
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

  &.primary-outline {
    border: 1px solid $c_gold;
    color: $c_gold;
    transition: all $trans_fast;
    &:hover {
      background: $c_gold;
      color: white;
    }
  }
  &.round {
    border-radius: 10px;
  }
  &.secondary {
    background: $c_gray_100;
    border: 1px solid $c_gray_200;
    color: $c_text;
    transition: background $trans_fast;

    &:hover {
      background: $c_gray_300;
    }
  }

  &.connect {
    background: $g_gray;
    width: 100%;
    padding: 24px 0;
    display: flex;
    justify-content: center;
    border: 1px solid $c_gold;
    border-radius: 24px;
    border-radius: 12px;
    &:hover {
      background: $g_gold;
      color: $c_white;
    }
    &.active {
      background: $g_gold;
      color: $c_white;
      &:hover {
        background: $g_gray;
        color: initial;
      }
    }
    &:disabled {
      background: $g_gold !important;
      color: $c_white !important;
    }
  }

  &.ghost {
    background: transparent;
    &.secondary {
      color: $c_gray_300;
      border: 2px solid $c_gray_300;
      &:hover {
        background: transparent;
        border: 2px solid $c_gray_400;
      }
    }
    &.success {
      color: #689829;
      border: 2px solid #689829;
      &:hover {
        background: transparent;
        border: 2px solid #689829;
      }
    }
    &.primary {
      color: $c_gold;
      border: 2px solid $c_gold;

      &:hover {
        background: transparent;
        border: 2px solid $c_gray_400;
      }
    }
  }

  // sizes:
  // block spans the full width of parent

  &.small {
    font-family: Arial, Helvetica, sans-serif;
    letter-spacing: 0.5px;
    font-style: normal;
    font-weight: bold;
    font-size: 10px;
    padding: 2px 5px;
    height: auto;
    line-height: initial;
    position: relative;
    top: -1px;
    &:active {
      transform: translateY(1px);
    }
  }

  &.block {
    display: block;
    width: 100%;
    margin: 0;
  }

  &.medium {
    padding: 0 30px;
  }
  & .content {
    display: flex;
    justify-content: center;
  }
}
</style>
