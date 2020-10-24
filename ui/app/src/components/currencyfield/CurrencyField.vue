<template>
  <div class="currency-field">
    <label class="label">{{ label }}</label>
    <label class="balance right-col"
      ><BalanceField :symbol="localBalance.symbol"
    /></label>
    <input
      class="input"
      type="number"
      v-model="localBalance.amount"
      @click="$event.target.select()"
    />

    <button @click="handleSelectClicked(localBalance)" class="button right-col">
      <span class="select-button" v-if="localBalance.symbol !== null">
        <AssetItem :symbol="localBalance.symbol" /><span>▾</span></span
      >
      <span v-else>Select</span>
    </button>
  </div>
</template>


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
    modelValue: Object,
  },
  components: { BalanceField, AssetItem },
  setup(props, context) {
    const { handleClicked: handleSelectClicked } = useSelectTokens();
    const localBalance = computed({
      get: () => props.modelValue,
      set: (balance) => context.emit("update:modelValue", balance),
    });

    return { handleSelectClicked, localBalance };
  },
});
</script>

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