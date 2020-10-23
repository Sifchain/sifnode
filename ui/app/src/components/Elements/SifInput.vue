<template>
  <div class="sif-input__wrapper">
    <label class="sif-input__label">{{ label }}</label>
    <input class="sif-input" type="text">
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: String
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
  border: none;

  &:focus {
    outline: none;
  }
}
</style>