


<script lang="ts">
import { defineComponent } from "vue";
import { computed } from "@vue/reactivity";
import BalanceField from "./BalanceField.vue";
import AssetItem from "@/components/tokenSelector/AssetItem.vue";
import SifButton from "@/components/shared/SifButton.vue";
import SifInput from "@/components/shared/SifInput.vue";

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
  emits: ["focus", "blur", "selectsymbol", "update:amount", "update:symbol"],
  components: { BalanceField, AssetItem, SifButton, SifInput },
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
    <div class="left">
      <label class="label">{{ label }}</label>
      <SifInput
        bold
        v-bind="$attrs"
        type="number"
        v-model="localAmount"
        @focus="$emit('focus', $event.target)"
        @blur="$emit('blur', $event.target)"
        @click="$event.target.select()"
      />
    </div>

    <div class="right">
      <label class="label">
        <BalanceField :symbol="localSymbol" />
      </label>
      <SifButton
        v-if="localSymbol !== null"
        secondary
        block
        @click="$emit('selectsymbol')"
      >
        <span><AssetItem :symbol="localSymbol" /> ▾ </span>
      </SifButton>
      <SifButton v-else primary block @click="$emit('selectsymbol')">
        <span>Select</span>
      </SifButton>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.currency-field {
  padding: 4px 15px 15px 15px;
  border-radius: $br_sm;
  background: $g_gray_reverse;
  color: $c_gray_700;
  display: flex;
}

.left,
.right {
  display: flex;
  flex-direction: column;
}

.left {
  align-items: flex-start;
  flex-grow: 1;
  margin-right: 10px;
}

.right {
  align-items: flex-end;
  width: 128px;
}

.label {
  font-size: $fs_sm;
}
</style>