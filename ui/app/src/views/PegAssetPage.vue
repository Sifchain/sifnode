<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { computed, ref, toRefs } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import {
  Asset,
  AssetAmount,
  Fraction,
  SwapState,
  useSwapCalculator,
} from "ui-core";
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
import RaisedPanel from "@/components/shared/RaisedPanel.vue";
import { useRouter } from "vue-router";
import SifInput from "@/components/shared/SifInput.vue";
import Label from "@/components/shared/Label.vue";
import RaisedPanelColumn from "@/components/shared/RaisedPanelColumn.vue";
import { trimZeros } from "ui-core/src/hooks/utils";
import BigNumber from "bignumber.js";

export default defineComponent({
  components: {
    Layout,
    CurrencyField,
    RaisedPanel,
    Label,
    SifInput,
    ActionsPanel,
    RaisedPanelColumn,
  },

  setup(props, context) {
    const { store, actions } = useCore();
    const router = useRouter();
    const mode = computed(() => {
      return router.currentRoute.value.path.indexOf("/peg/reverse") > -1
        ? "unpeg"
        : "peg";
    });

    // const symbol = ref<string | null>(null);
    const symbol = computed(() => {
      const assetFrom = router.currentRoute.value.params.assetFrom;
      return Array.isArray(assetFrom) ? assetFrom[0] : assetFrom;
    });
    const amount = ref("0.0");
    const address = ref(store.wallet.sif.address);

    async function handlePeg() {
      await actions.peg.lock(
        address.value,
        AssetAmount(Asset.get(symbol.value), amount.value)
      );
    }

    async function handleUnpeg() {
      await actions.peg.burn(
        address.value,
        AssetAmount(Asset.get(symbol.value), amount.value)
      );
    }
    return {
      mode,
      symbol,
      amount,
      address,
      handleBlur: () => {
        amount.value = trimZeros(amount.value);
      },
      handleSelectSymbol: () => {},
      handleMaxClicked: () => {
        const balances =
          mode.value === "peg"
            ? store.wallet.eth.balances
            : store.wallet.sif.balances;
        const accountBalance = balances.find((balance) => {
          return (
            balance.asset.symbol.toLowerCase() === symbol.value.toLowerCase()
          );
        });

        if (!accountBalance) return;

        amount.value = accountBalance.toFixed();
      },
      handleAmountUpdated: (newAmount: string) => {
        amount.value = newAmount;
      },
      handleActionClicked: () => {
        if (mode.value === "peg") {
          handlePeg();
        } else {
          handleUnpeg();
        }
      },

      nextStepAllowed: computed(() => {
        const amountNum = new BigNumber(amount.value);
        return amountNum.isGreaterThan("0.0") && address.value !== "";
      }),
    };
  },
});
</script>

<template>
  <Layout :title="mode === 'peg' ? 'Peg Asset' : 'Unpeg Asset'" backLink="/peg">
    <div class="vspace">
      <CurrencyField
        :amount="amount"
        :max="true"
        :selectable="true"
        :symbol="symbol"
        :symbolFixed="true"
        @blur="handleBlur"
        @maxclicked="handleMaxClicked"
        @update:amount="handleAmountUpdated"
        label="From"
      />
      <RaisedPanel>
        <RaisedPanelColumn v-if="mode === 'peg'">
          <Label>Sifchain Recipient Address</Label>
          <SifInput
            v-model="address"
            placeholder="Eg. sif21syavy2npfyt9tcncdtsdzf7kny9lh777yqcnd"
          />
        </RaisedPanelColumn>
        <RaisedPanelColumn v-if="mode === 'unpeg'">
          <Label>Ethereum Recipient Address</Label>
          <SifInput
            v-model="address"
            placeholder="Eg. 0xeaf65652e380528fffbb9fc276dd8ef608931e3c"
          />
        </RaisedPanelColumn>
      </RaisedPanel>
      <ActionsPanel
        connectType="connectToAll"
        @nextstepclick="handleActionClicked"
        :nextStepAllowed="nextStepAllowed"
        :nextStepMessage="mode === 'peg' ? 'Peg' : 'Unpeg'"
      />
    </div>
  </Layout>
</template>

<style lang="scss" scoped>
.vspace {
  display: flex;
  flex-direction: column;
  & > * {
    margin-bottom: 1rem;
  }

  & > *:last-child {
    margin-bottom: 0;
  }
}
</style>