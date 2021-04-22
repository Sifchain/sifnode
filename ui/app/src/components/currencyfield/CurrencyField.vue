<script lang="ts">
import { defineComponent } from "vue";
import { computed } from "@vue/reactivity";
import BalanceField from "./BalanceField.vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import SifButton from "@/components/shared/SifButton.vue";
import SifInput from "@/components/shared/SifInput.vue";
import Caret from "@/components/shared/Caret.vue";
import RaisedPanel from "@/components/shared/RaisedPanel.vue";
import Label from "@/components/shared/Label.vue";
import RaisedPanelColumn from "@/components/shared/RaisedPanelColumn.vue";
import RaisedPanelRow from "@/components/shared/RaisedPanelRow.vue";
import Checkbox from "@/components/shared/Checkbox.vue";

export type BalanceShape = {
  symbol: string;
  amount: string;
  available: string;
};

export default defineComponent({
  props: {
    label: String,
    slug: String,
    amount: String,
    symbol: String,
    available: String,
    inputDisabled: { type: Boolean, default: false },
    selectable: { type: Boolean, default: true },
    max: { type: Boolean, default: false },
    isMaxActive: { type: Boolean, default: false },
    symbolFixed: { type: Boolean, default: false },
    toggleLabel: { type: String, default: null },
    asyncPooling: { type: Boolean, default: null },
    handleToggle: { type: Function, default: null },
  },
  inheritAttrs: false,
  emits: [
    "focus",
    "blur",
    "selectsymbol",
    "update:amount",
    "update:symbol",
    "handleToggle",
    "maxclicked",
  ],
  components: {
    RaisedPanelColumn,
    RaisedPanelRow,
    RaisedPanel,
    BalanceField,
    AssetItem,
    SifButton,
    Caret,
    SifInput,
    Label,
    Checkbox,
  },
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
  methods: {
    onClickChild(value: string) {
      this.handleToggle();
    },
  },
});
</script>

<template>
  <RaisedPanel>
    <Checkbox
      @clicked="onClickChild"
      v-if="toggleLabel"
      :toggleLabel="toggleLabel"
      :checked="asyncPooling"
    />
    <RaisedPanelRow>
      <RaisedPanelColumn class="left">
        <Label>{{ label }}</Label>
        <SifInput
          bold
          :data-handle="slug + '-input'"
          v-bind="$attrs"
          type="number"
          v-model="localAmount"
          :disabled="inputDisabled"
          @focus="$emit('focus', $event.target)"
          @blur="$emit('blur', $event.target)"
          ><template v-slot:end
            ><SifButton
              v-if="max && !isMaxActive"
              :data-handle="slug + '-max-button'"
              @click="$emit('maxclicked')"
              small
              ghost
              primary
              >MAX</SifButton
            ></template
          ></SifInput
        >
      </RaisedPanelColumn>
      <RaisedPanelColumn class="right">
        <Label>
          <BalanceField :symbol="localSymbol" />
        </Label>

        <SifButton
          :data-handle="slug + '-select-button'"
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
          :data-handle="slug + '-select-button'"
          :disabled="!selectable"
          v-if="localSymbol === null"
          primary
          block
          @click="$emit('selectsymbol')"
        >
          <span>Select</span>
        </SifButton>
      </RaisedPanelColumn>
    </RaisedPanelRow>
  </RaisedPanel>
</template>

<style lang="scss" scoped>
.fixed-symbol {
  width: 100%;
  margin: auto;
  display: flex;
  justify-content: center;
  pointer-events: none;
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
</style>
