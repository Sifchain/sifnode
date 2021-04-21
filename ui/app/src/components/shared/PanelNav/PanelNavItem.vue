<script lang="ts">
import { defineComponent } from "vue";
import Icon from "../Icon.vue";
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
        gold: this.color === "gold",
      },
    };
  },
  methods: {
    subPageIsActive(input: string) {
      const paths = [input];
      return paths.some((path) => {
        return this.$route?.path.indexOf(path) === 0; // current path starts with this path string
      });
    },
  },
});
// :class="{ disabled: disabled, 'router-link-active': subPageIsActive(to) }"
</script>

<template>
  <router-link
    class="link"
    :to="to"
    :class="{ disabled: disabled, 'router-link-active': false }"
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
  border-radius: 10px;
  border: 2px solid $c_gray_700;
  overflow: hidden;
  transition: border-color $trans_fast;
  cursor: pointer;
  background: linear-gradient(180deg, #5b5b5b 0%, #4d4d4d 0.01%, #262626 100%);
}

.label {
  color: $c_gray_400;
  transition: color $trans_fast;
  text-align: center;
  margin-top: 5px;
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
  &:hover {
    .icon {
      border-color: #525252;
    }

    .label {
      color: $c_gray_200;
    }
  }
  .icon-svg {
    opacity: 0.5;
  }
  &.router-link-active {
    .icon {
      background: linear-gradient(180deg, #caa93a 0%, #a6820b 100%);
      pointer-events: none;
      border-color: #68530d;
    }
    &:hover {
      .icon {
        border-color: #68530d;
      }

      .label {
        color: $c_gray_200;
      }
    }
  }
}
</style>
