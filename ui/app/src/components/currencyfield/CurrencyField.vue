


<script lang="ts">
import { defineComponent } from "vue";
import { computed } from "@vue/reactivity";
import Modal from "@/components/shared/Modal.vue";
import BalanceField from "./BalanceField.vue";
import AssetItem from "@/components/tokenSelector/AssetItem.vue";

export type BalanceShape = {
  symbol: string;
  amount: string;
  available: string;
};

export default defineComponent({
  props: {
    label: String,
    amount: String,
    symbol: String,
    available: String,
  },
  inheritAttrs: false,
  emits: ["selectsymbol", "update:amount", "update:symbol"],
  components: { BalanceField, AssetItem, Modal },
  setup(props, context) {
    const localAmount = computed({
      get: () => props.amount,
      set: (amount) => context.emit("update:amount", amount),
    });

    const localSymbol = computed({
      get: () => props.symbol,
      set: (symbol) => context.emit("update:symbol", symbol),
    });

    return { localSymbol, localAmount };
  },
});
</script>

<template>
  <div class="currency-field">
    <label class="label">{{ label }}</label>
    <label class="balance right-col"
      ><BalanceField :symbol="localSymbol"
    /></label>
    <input
      v-bind="$attrs"
      class="input"
      type="number"
      v-model="localAmount"
      @focus="$emit('focus', $event.target)"
      @blur="$emit('blur', $event.target)"
      @click="$event.target.select()"
    />

    <button @click="$emit('selectsymbol')" class="button right-col">
      <span class="select-button" v-if="localSymbol !== null">
        <AssetItem :symbol="localSymbol" /><span>▾</span></span
      >
      <span v-else>Select</span>
    </button>
  </div>
</template>

<style scoped>
.currency-field {
  border: 1px solid grey;
  padding: 1rem;
  display: grid;
  grid-gap: 1rem;
  grid-template-areas: "label balance" "input button";
}
.label {
  grid-area: "label";
}
.right-col {
  width: 6rem;
}
.balance {
  grid-area: "balance";
}

.input {
  grid-area: "input";
}
.button {
  grid-area: "button";
}
.select-button {
  display: flex;
  flex-direction: row;
  align-items: center;
}
</style>