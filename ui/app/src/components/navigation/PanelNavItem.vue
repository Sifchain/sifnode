<template>
  <router-link
    class="link"
    :to="to"
    :class="{ disabled: disabled, 'router-link-active': subPageIsActive(to) }"
  >
    <div class="icon-button">
      <div class="icon" :class="colorClass"></div>
      <Icon class="icon-svg" :icon="icon" />
    </div>
    <div class="label" v-if="label" :class="colorClass">
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
    color: {
      type: String,
      default: "blue",
    },
  },
  data() {
    return {
      colorClass: {
        green: this.color === "green",
        blue: this.color === "blue",
        pink: this.color === "pink",
        gold: this.color === "gold"
      },
    };
  },
  methods: {
    subPageIsActive(input: string) {
      const paths = [input];
      return paths.some((path) => {
        return this.$route.path.indexOf(path) === 0; // current path starts with this path string
      });
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
      &.blue {
        color: $c_blue;
      }
      &.green {
        color: $c_green;
      }
    }
    .icon-svg {
      opacity: 1;
    }
    .icon.blue {
      border-color: $c_blue;
      &::after {
        background: linear-gradient(to bottom, #61a5f6 0%, #a0caf9 100%);
      }
      &::before {
        background: linear-gradient(to bottom, #8bbef8 0%, #90c2f9 100%);
      }
    }
    .icon.green {
      border-color: $c_green;
      &::after {
        background: linear-gradient(to bottom, #a0dd2b 0%, #a5de30 100%);
      }
      &::before {
        background: linear-gradient(to bottom, #8fd817 0%, #b1e245 100%);
      }
    }
    .icon.pink {
      border-color: $c_pink;
      &::after {
        background: linear-gradient(to bottom, #f09680 0%, #f1a3ac 100%);
      }
      &::before {
        background: linear-gradient(to bottom, #f3afb7 0%, #f1a4ae 100%);
      }
    }
    .icon.gold {
      border-color: $c_gold;
      &::after {
        background: linear-gradient(to bottom, $c_gold 0%, #fbe59d 100%);
      }
      &::before {
        background: linear-gradient(to bottom, #fbe59d 0%, $c_gold 100%);
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
