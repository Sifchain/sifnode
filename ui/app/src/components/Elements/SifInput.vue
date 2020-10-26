<template>
  <div class="sif-input__wrapper">
    <label class="sif-input__label">{{ label }}</label>
    <input class="sif-input" type="text" :class="classes" :placeholder="placeholder">
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: String
    },
    gold: {
      type: Boolean,
      default: false,
    },
    placeholder: {
      type: String,
      default: '',
    }
  },

  data() {
    return {
      classes: {
        'sif-input--gold': this.gold,
      }
    }
  }
}
</script>

<style lang="scss">
.sif-input {
  &__wrapper {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    width: 100%;
  }

  height: 30px;
  width: 100%;
  padding: 0 8px;
  border-radius: $br_sm;
  border: 1px solid $c_white;

  &::placeholder {
    font-family: $f_default;
    color: $c_gray_400;
  }

  &:focus {
    outline: none;
  }

  &--gold {
    border-color: $c_gold;
  }
}
</style>