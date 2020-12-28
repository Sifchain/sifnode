<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { computed, ref, toRefs } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { Asset, SwapState, useSwapCalculator } from "ui-core";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import CurrencyField from "@/components/currencyfield/CurrencyField.vue";
import Modal from "@/components/shared/Modal.vue";
import SelectTokenDialogEth from "@/components/tokenSelector/SelectTokenDialogEth.vue";
import SelectTokenDialogSif from "@/components/tokenSelector/SelectTokenDialogSif.vue";
import PriceCalculation from "@/components/shared/PriceCalculation.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationDialog, {
  ConfirmState,
} from "@/components/confirmationDialog/ConfirmationDialog.vue";
import { useCurrencyFieldState } from "@/hooks/useCurrencyFieldState";
import DetailsPanel from "@/components/shared/DetailsPanel.vue";

import { useRouter } from "vue-router";

export default defineComponent({
  components: {
    Modal,
    Layout,
    SelectTokenDialogSif,
    CurrencyField,
    SelectTokenDialogEth,
  },

  setup(_, context) {
    const router = useRouter();
    const mode = computed(() => {
      return router.currentRoute.value.path.indexOf("unpeg") > -1
        ? "unpeg"
        : "peg";
    });

    const symbol = ref<string | null>(null);

    return {
      mode,
      symbol,
      amount: ref("0"),
      handleBlur: () => {},
      handleSelectSymbol: () => {},
      handleMaxClicked: () => {},
      handleUpdateAmount: () => {},
      handleFromUpdateSymbol: () => {},
      handleSelectClosed(data: string) {
        if (typeof data !== "string") {
          return;
        }

        symbol.value = data;
      },
    };
  },
});
</script>

<template>
  <Layout backLink="/peg">
    <Modal @close="handleSelectClosed">
      <template v-slot:activator="{ requestOpen }">
        <CurrencyField
          :amount="amount"
          :max="true"
          :selectable="true"
          :symbol="symbol"
          :symbolFixed="false"
          @blur="handleBlur"
          @maxclicked="handleMaxClicked"
          @selectsymbol="requestOpen"
          @update:amount="handleUpdateAmount"
          @update:symbol="handleFromUpdateSymbol"
          label="From"
        />
      </template>
      <template v-slot:default="{ requestClose }">
        <SelectTokenDialogEth
          v-if="mode === 'peg'"
          :selectedTokens="[symbol]"
          @tokenselected="requestClose"
        />
        <SelectTokenDialogSif
          v-if="mode === 'unpeg'"
          :selectedTokens="[symbol]"
          @tokenselected="requestClose"
        />
      </template>
    </Modal>
  </Layout>
</template>