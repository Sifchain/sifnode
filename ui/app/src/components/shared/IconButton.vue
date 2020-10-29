<template>
  <div class="wrapper">
    <div class="icon-button" :class="classes">

    </div>
    <div class="label" v-if="label">
      {{ label }}
    </div>
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: String
    },
    active: {
      type: Boolean,
      default: false
    },
    disabled: {
      type: Boolean,
      default: false
    },
  },

  data() {
    return {
      classes: {
        'disabled': this.disabled,
        'active': this.active,
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.wrapper {
  width: 55px;
}
.icon-button {
  position: relative;
  width: 55px;
  height: 55px;
  display: block;
  border-radius: $br_md;
  border: 2px solid $c_gray_700;
  overflow: hidden;
  transition: border-color $trans_fast;
  cursor: pointer;

  &::after {
    content: "";
    display: block;
    width: 75px;
    height: 75px;
    position: absolute;
    background: linear-gradient(to bottom, #3D3D3D 0%, #7A7A7A 100%);
    transform: rotate(45deg);
    top: 10px;
    left: 17px;
  }
  &::before {
    content: "";
    display: block;
    width: 55px;
    height: 55px;
    position: absolute;
    background: linear-gradient(to bottom, #707070 0%, #575757 90%);
  }

  .label {
    color: $c_gray_400;
    transition: color $trans_fast;
  }

  &:hover {
    border-color: $c_gray_200;

    & ~ .icon-button__label {
      color: $c_gray_200;
    }
  }

  &.disabled {
    pointer-events: none;
    border-color: transparent;
    opacity: 0.38;

    & ~ .label {
      color: rgba($c_gray_700, 0.38);
    }
  }

  &.active {
    border-color: $c_blue;
    pointer-events: none;

    & ~ .label {
      color: $c_white;
    }

    &::after {
      background: linear-gradient(to bottom, #61A5F6 0%, #A0CAF9 100%);
    }

    &::before {
      background: linear-gradient(to bottom, #8BBEF8 0%, #90C2F9 100%);
    }
  }
}
</style>