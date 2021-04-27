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
    center: {
      type: Boolean,
      default: false,
    },
    icon: {
      type: String,
    },
  },
  setup(props) {
    const { center } = props;
    const iconClasses = {
      "icon-svg": true,
      "icon-svg-center": center,
    };
    return { iconClasses };
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
</script>

<template>
  <router-link
    class="link"
    :to="to"
    :class="{ disabled: disabled, 'router-link-active': false }"
  >
    <div class="icon-button">
      <div class="icon"></div>
      <Icon :class="iconClasses" :icon="icon" />
    </div>
    <div class="label" v-if="label">
      {{ label }}
    </div>
  </router-link>
</template>

<style lang="scss" scoped>
.link {
  text-decoration: none;
  font-size: 13px;
  width: 65px;
  position: relative;
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

    // background: linear-gradient(360deg, #f2d267 0%, red 100%);
    background: linear-gradient(180deg, #525252 0%, #151515 90%, #000000 100%);
    border-radius: 12px;
    width: 58px;
    height: 58px;
  }
  .icon {
    position: relative;
    width: 55px;
    height: 55px;
    display: block;
    border-radius: 10px;
    overflow: hidden;
    transition: border-color $trans_fast;
    cursor: pointer;
    // background: linear-gradient(180deg, #5b5b5b 0%, #4d4d4d 0.01%, blue 100%);
    background: linear-gradient(
      180deg,
      #5b5b5b 0%,
      #4d4d4d 0.01%,
      #262626 100%
    );
  }

  .label {
    color: $c_gray_400;
    transition: color $trans_fast;
    text-align: center;
    margin-top: 5px;
  }

  &.disabled {
    pointer-events: none;
    opacity: 0.38;

    .label {
      color: rgba($c_gray_700, 0.38);
    }
  }
  &:hover {
    .label {
      color: $c_gray_200;
    }
  }
  .icon-svg {
    opacity: 0.5;
  }
  .icon-svg-center {
    display: flex;
    align-items: center;
  }
  &.router-link-active {
    .icon-button {
      background: linear-gradient(180deg, #f2d267 0%, #68530d 100%);
    }
    .icon {
      background: linear-gradient(180deg, #caa93a 0%, #a6820b 100%);
      pointer-events: none;
    }
    &:hover {
      .label {
        color: $c_gray_200;
      }
    }
  }
}
</style>
