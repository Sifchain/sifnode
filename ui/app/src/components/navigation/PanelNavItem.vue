<template>
  <router-link class="link" :to="to" :class="{ disabled: disabled }">
    <div class="icon-button">
      <div class="icon"></div>
      <Icon class="icon-svg" :icon="icon" />
    </div>
    <div class="label" v-if="label">
      {{ label }}
    </div>
  </router-link>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Icon from "@/components/shared/Icon.vue";
export default defineComponent({
  components: { Icon },
  props: {
    to: {
      type: String,
    },
    label: {
      type: String,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    icon: {
      type: String,
    },
  },
});
</script>

<style lang="scss" scoped>
.icon-svg {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  height: 66px;
  stroke: white;
}
.icon-button {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}
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
    background: linear-gradient(to bottom, #3d3d3d 0%, #7a7a7a 100%);
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
  .icon-svg {
    opacity: 0.5;
  }
  &.router-link-active {
    pointer-events: none;

    .label {
      color: $c_white;
    }
    .icon-svg {
      opacity: 1;
    }
    .icon {
      border-color: $c_blue;
      &::after {
        background: linear-gradient(to bottom, #61a5f6 0%, #a0caf9 100%);
      }
      &::before {
        background: linear-gradient(to bottom, #8bbef8 0%, #90c2f9 100%);
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