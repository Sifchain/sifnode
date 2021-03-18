<template>
  <div class="slider">
    <p class="text--left" v-if="message">{{ message }}</p>
    <input
      :disabled="disabled"
      v-model="localValue"
      class="input"
      :min="min"
      :max="max"
      type="range"
      :step="step"
    />
    <div class="row">
      <div>
        <label class="label" @click="$emit('leftclicked')">{{
          leftLabel
        }}</label>
      </div>
      <div>
        <label class="label" @click="$emit('middleclicked')">{{
          middleLabel
        }}</label>
      </div>
      <div>
        <label class="label" @click="$emit('rightclicked')">{{
          rightLabel
        }}</label>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { computed } from "@vue/reactivity";
export default defineComponent({
  emits: ["leftclicked", "middleclicked", "rightclicked", "update:modelValue"],
  props: {
    message: { type: String, default: "" },
    disabled: { type: Boolean, default: false },
    modelValue: { type: String, default: "0" },
    min: { type: String, default: "0" },
    max: { type: String, default: "100" },
    step: { type: String, default: "1" },
    leftLabel: { type: String, default: "0" },
    middleLabel: { type: String, default: "50" },
    rightLabel: { type: String, default: "100" },
  },
  setup(props, context) {
    const localValue = computed({
      get: () => props.modelValue,
      set: (value) => context.emit("update:modelValue", value),
    });
    return { localValue };
  },
});
</script>

<style lang="scss">
.slider {
  margin-bottom: 1rem;
  width: 100%;
  .input {
    width: 100%;
    cursor: pointer;
  }
  .row {
    display: flex;
    justify-content: space-between;
    & > * {
      width: 32%;
    }
    & > *:first-child {
      text-align: left;
    }
    & > *:last-child {
      text-align: right;
    }

    .label {
      transition: opacity $trans_fast;
      &:hover {
        cursor: pointer;
        color: $c_black;
      }
    }
  }
}
</style>
