<template>
  <div class="sif-input-wrapper">
    <label class="label">{{ label }}</label>
    <div class="sif-input-background" :class="classes">
      <slot name="start"></slot>
      <input
        v-bind="$attrs"
        v-model="localValue"
        class="sif-input-control"
        :disabled="disabled"
        :placeholder="placeholder"
      />
      <slot name="end"></slot>
    </div>
  </div>
</template>

<script>
import { computed, defineComponent } from "vue";

export default defineComponent({
  inheritAttrs: false,

  props: {
    label: {
      type: String,
    },
    modelValue: { type: String },
    gold: {
      type: Boolean,
      default: false,
    },
    placeholder: {
      type: String,
      default: "",
    },
    bold: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["update:modelValue"],
  setup(props, context) {
    const localValue = computed({
      get: () => props.modelValue,
      set: (value) => context.emit("update:modelValue", value),
    });
    return {
      localValue,
      classes: {
        gold: props.gold,
        bold: props.bold,
      },
    };
  },
});
</script>

<style lang="scss" scoped>
.sif-input-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  width: 100%;
}

.sif-input-background {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: stretch;
  height: 30px;
  width: 100%;
  background: $c_white;
  box-sizing: border-box;
  border-radius: $br_sm;
  border: 1px solid $c_gray_200;
  padding: 0 3px;

  &.gold {
    border-color: $c_gold;
  }

  &.bold {
    font-weight: bold;
  }
}

.sif-input-control {
  padding: 0 8px;
  border: none;
  box-sizing: border-box;
  flex: 1;
  display: block;
  width: 100%;

  &::placeholder {
    font-family: $f_default;
    color: $c_gray_400;
  }

  &:focus {
    outline: none;
  }
}
</style>
