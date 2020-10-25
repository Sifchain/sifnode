


<script lang="ts">
import { defineComponent } from "vue";
import { computed } from "@vue/reactivity";
import { useSelectTokens } from "./useSelectToken";
import BalanceField from "./BalanceField.vue";
import AssetItem from "./AssetItem.vue";

export type BalanceShape = {
  symbol: string;
  amount: string;
  available: string;
};

export default defineComponent({
  props: {
    label: String,
    modelkey: String,
    amount: String,
    symbol: String,
    available: String,
  },
  components: { BalanceField, AssetItem },
  setup(props, context) {
    const { handleClicked: handleSelectClicked } = useSelectTokens();
    const localAmount = computed({
      get: () => props.amount,
      set: (amount) => context.emit("update:amount", amount),
    });

    const localSymbol = computed({
      get: () => props.symbol,
      set: (symbol) => context.emit("update:symbol", symbol),
    });

    return { handleSelectClicked, localSymbol, localAmount };
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
      class="input"
      type="number"
      v-model="localAmount"
      @focus="$emit('focus', $event.target)"
      @click="$event.target.select()"
    />

    <button @click="handleSelectClicked(modelkey)" class="button right-col">
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