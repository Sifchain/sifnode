<template>
  <router-link class="link" :to="to" :class="{'disabled': disabled}">
    <div class="icon"></div>
    <div class="label" v-if="label">
      {{ label }}
    </div>
  </router-link>
</template>

<script>
export default {
  props: {
    to: {
      type: String
    },
    label: {
      type: String
    },
    disabled: {
      type: Boolean,
      default: false
    },
  },
}
</script>

<style lang="scss" scoped>
.icon {
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
}

.label {
  color: $c_gray_400;
  transition: color $trans_fast;
}

.link {
  width: 55px;
  text-decoration: none;

  &.disabled {
    pointer-events: none;
    border-color: transparent;
    opacity: 0.38;

    .label {
      color: rgba($c_gray_700, 0.38);
    }
  }

  &.router-link-active {
    pointer-events: none;

    .label {
      color: $c_white;
    }

    .icon {
      border-color: $c_blue;
      &::after {
        background: linear-gradient(to bottom, #61A5F6 0%, #A0CAF9 100%);
      }
      &::before {
        background: linear-gradient(to bottom, #8BBEF8 0%, #90C2F9 100%);
      }
    }
  }
  
  &:hover {
    .icon {
      border-color: $c_gray_200;
    }

    .label {
      color: $c_gray_200;
    }
  }
}

</style>