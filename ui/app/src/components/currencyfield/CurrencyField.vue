


<script lang="ts">
import { defineComponent } from "vue";
import { computed } from "@vue/reactivity";
import BalanceField from "./BalanceField.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import SifButton from "@/components/shared/SifButton.vue";
import SifInput from "@/components/shared/SifInput.vue";
import Caret from "@/components/shared/Caret.vue";

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
    selectable: { type: Boolean, default: true },
    max: { type: Boolean, default: false },
    symbolFixed: { type: Boolean, default: false },
  },
  inheritAttrs: false,
  emits: [
    "focus",
    "blur",
    "selectsymbol",
    "update:amount",
    "update:symbol",
    "maxclicked",
  ],
  components: { BalanceField, AssetItem, SifButton, Caret, SifInput },
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
        ><template v-slot:end
          ><SifButton
            v-if="max"
            class="max-button"
            @click="$emit('maxclicked')"
            small
            ghost
            >MAX</SifButton
          ></template
        ></SifInput
      >
    </div>

    <div class="right">
      <label class="label">
        <BalanceField :symbol="localSymbol" />
      </label>

      <SifButton
        nocase
        v-if="localSymbol !== null && !symbolFixed"
        secondary
        block
        @click="$emit('selectsymbol')"
      >
        <span><AssetItem :symbol="localSymbol" /></span>
        <span><Caret /></span>
      </SifButton>
      <div v-if="localSymbol !== null && symbolFixed">
        <AssetItem :symbol="localSymbol" />
      </div>

      <SifButton
        :disabled="!selectable"
        v-if="localSymbol === null"
        primary
        block
        @click="$emit('selectsymbol')"
      >
        <span>Select</span>
      </SifButton>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.fixed-symbol {
  width: 100%;
  margin: auto;
  display: flex;
  justify-content: center;
  pointer-events: none;
}
.currency-field {
  padding: 4px 15px 15px 15px;
  border-radius: $br_sm;
  border: 1px solid $c_gray_100;
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
  flex-shrink: 0;
}

.label {
  font-size: $fs_sm;
  white-space: nowrap;
}
</style>