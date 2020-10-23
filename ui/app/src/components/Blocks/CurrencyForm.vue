<template>
   <div class="currency-form text--small">
      <div class="currency-form__input">
        <SifInput :label="label" />
      </div>
      <div class="currency-form__balance">
        <span>Balance: 1200,7356</span>
        <SifButton primary medium>Select</SifButton>
      </div>
    </div>
</template>

<script>
import SifButton from '../Elements/SifButton.vue';
import SifInput from '../Elements/SifInput.vue';

export default {
  props: {
    label: {
      type: String,
      default: ''
    }
  },

  components: {
    SifButton,
    SifInput,
  }
}
</script>

<style lang="scss">
.currency-form {
  display: flex;

  &__input,
  &__balance {
    display: flex;
    flex-direction: column;
  }

  &__balance {
    align-items: flex-end;
  }

  &__input {
    align-items: flex-start;
    flex-grow: 1;
    margin-right: 10px;
  }
}
</style>