<template>
  <div class="sif-input-wrapper">
    <label class="label">{{ label }}</label>
    <input
      v-bind="$attrs"
      v-model="localValue"
      class="sif-input-control"
      :class="classes"
      :placeholder="placeholder"
    />
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

.sif-input-control {
  height: 30px;
  width: 100%;
  padding: 0 8px;
  box-sizing: border-box;
  border-radius: $br_sm;
  border: 1px solid $c_white;

  &::placeholder {
    font-family: $f_default;
    color: $c_gray_400;
  }

  &:focus {
    outline: none;
  }

  &.gold {
    border-color: $c_gold;
  }

  &.bold {
    font-weight: bold;
  }
}
</style>