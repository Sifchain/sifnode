

<template>
  <div class="slippage-section mh-70">
    <div class="d-flex fd-row">
      <p class="text--small text--italic text--normal" align="left">
        Slippage tolerance
      </p>
    </div>
    <div class="field-wrappers d-flex fd-row j-start">
      <div class="slippage-item" v-bind:class="{ 'slippage-item__selected': localSlippage === '0.5' }" @click="handleUpdateSlippage('0.5')">
        <p class="text--standard text--normal text--regular">
          0.5%
        </p>
      </div>
      <div class="slippage-item" v-bind:class="{ 'slippage-item__selected': localSlippage === '1.0' }" @click="handleUpdateSlippage('1.0')">
        <p class="text--standard text--normal text--regular">
          1.0%
        </p>
      </div>
      <div class="slippage-item" v-bind:class="{ 'slippage-item__selected': localSlippage === '2.0' }" @click="handleUpdateSlippage('2.0')">
        <p class="text--standard text--normal text--regular">
          2.0%
        </p>
      </div>
      <div class="sif-input__wrapper custom-slippage-input-container">
        <input
          @click="$event.target.select()"
          @input="validateSlippage"
          v-model="localSlippage"
          class="sif-input num_percent"
          type="number"
          step="any"
          min="0"
        >
        <span>%</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { computed } from "@vue/reactivity";
export default defineComponent({
  components: { },
  props: {
    slippage: { type: String, default: "0" },
  },
  emits: [
    "update:slippage",
  ],
  setup(props, context) {
    const localSlippage = computed({
      get: () => props.slippage,
      set: (amount) => handleUpdateSlippage(amount),
    });

    function handleUpdateSlippage(amount: string) {
      context.emit("update:slippage", amount);
    }

    function isNumeric(s: any) {
      return (s - 0) == s && (''+s).trim().length > 0;
    }

    function validateSlippage() {
      try {
        if (!isNumeric(localSlippage.value)) {
          localSlippage.value = "1.0";
        } else if (parseFloat(localSlippage.value) < 0) {
          localSlippage.value = localSlippage.value.substring(1);
        } else if (parseFloat(localSlippage.value) >= 99.0) {
          localSlippage.value = "99.0"
        } else if (!(localSlippage.value.startsWith(".") && localSlippage.value.length === 2)) {
          localSlippage.value = parseFloat(localSlippage.value).toFixed(1);
        }
      } catch (e) {
        localSlippage.value = "1.0";
      }
    }

    return {
      localSlippage,
      handleUpdateSlippage,
      validateSlippage
    };
  },
});
</script>

<style scoped lang="scss">
  .num_percent {
    margin: 5px;
    text-align: right;
    padding-right: 24px;
    font-size: 14px;
    font-weight: normal;
    font-style: normal;
    color: $c_text;
  }
  .sif-input {
    &__wrapper {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      width: 100%;
    }

    border-radius: $br_sm;
    border: 1px solid $c_gray_200;
    height: 30px;
    width: 100%;
    font-size: 15px;
    box-sizing: border-box;
    flex: 1;
    display: block;

    &::placeholder {
      font-family: $f_default;
      color: $c_gray_400;
    }

    &:focus {
      outline: none;
      border-color: $c_gold;
      border-width: 1px;
      border-style: solid;
    }
  }
  .slippage-section {
    margin-bottom: 20px;
    padding: 8px 4px 12px 12px;
    border-radius: $br_sm;
    border: 1px solid $c_gray_100;
    background: $g_gray_reverse;
    color: $c_gray_700;
  }
  .slippage-item {
    height: 30px;
    width: 55px;
    padding: 8px;
    margin: 5px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    border-radius: $br_sm;
    border: 1px solid $c_gray_400;
    background: $c_white;
    color: $c_text;

    &:hover {
      cursor: pointer;
    }

    &__selected {
      background: $c_gold;
      color: $c_white;
      border: none;
    }
  }
  .custom-slippage-input-container { position:relative; margin: 0 16px; }
  .custom-slippage-input-container span {
    position: absolute;
    right: 4px;
    top: 9px;
    font-family: $f_default;
    color: $c_text;
  }
</style>
