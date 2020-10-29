<template>
   <div class="currency-form text--small">
      <div class="input">
        <SifInput :label="label" />
      </div>
      <div class="balance">
        <span>Balance: 1200,7356</span>
        <SifSelect :tokens="tokens" @change="onChange" />
      </div>
    </div>
</template>

<script>
import SifButton from './SifButton.vue';
import SifInput from './SifInput.vue';
import SifSelect from './SifSelect.vue';

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
    SifSelect,
  },

  data() {
    return {
      tokens: [
        {
          icon: '',
          name: 'ETH',
        },
        {
          icon: '',
          name: '0xBTC',
        },
        {
          icon: '',
          name: 'aDAI',
        },
        {
          icon: '',
          name: 'AMN',
        },
      ]
    }
  },

  methods: {
    onChange(event, data) {
      console.log(data);
    }
  }
}
</script>

<style lang="scss" scoped>
.currency-form {
  display: flex;

  .input,
  .balance {
    display: flex;
    flex-direction: column;
  }

  .balance {
    align-items: flex-end;
  }

  .input {
    align-items: flex-start;
    flex-grow: 1;
    margin-right: 10px;
  }
}
</style>