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
import DetailsTable from "@/components/shared/DetailsTable.vue";
import Label from "@/components/shared/Label.vue";
import RaisedPanelColumn from "@/components/shared/RaisedPanelColumn.vue";
import { trimZeros } from "ui-core/src/hooks/utils";
import BigNumber from "bignumber.js";
import { useAssetItem } from "../components/shared/utils";
import B from "ui-core/src/entities/utils/B";

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export default defineComponent({
  components: {
    Layout,
    CurrencyField,
    RaisedPanel,
    Label,
    SifInput,
    DetailsTable,
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
    const address = computed(() =>
      mode.value === "peg" ? store.wallet.sif.address : store.wallet.eth.address
    );

    async function handlePeg() {
      try {
        await actions.peg.peg(
          AssetAmount(Asset.get(symbol.value), amount.value)
        );
        router.push("/peg");
      } catch (err) {
        console.error(err);
      }
    }

    async function handleUnpeg() {
      try {
        await actions.peg.unpeg(
          AssetAmount(Asset.get(symbol.value), amount.value)
        );
        router.push("/peg");
      } catch (err) {
        console.error(err);
      }
    }

    const accountBalance = computed(() => {
      const balances =
        mode.value === "peg"
          ? store.wallet.eth.balances
          : store.wallet.sif.balances;
      return balances.find((balance) => {
        return (
          balance.asset.symbol.toLowerCase() === symbol.value.toLowerCase()
        );
      });
    });

    const nextStepAllowed = computed(() => {
      const amountNum = new BigNumber(amount.value);
      const balance = accountBalance.value?.toFixed(0) ?? "0.0";
      return (
        amountNum.isGreaterThan("0.0") &&
        address.value !== "" &&
        amountNum.isLessThanOrEqualTo(balance)
      );
    });

    return {
      mode,
      modeLabel: computed(() => capitalize(mode.value)),
      symbol,
      symbolLabel: useAssetItem(symbol).label,
      amount,
      address,
      feeAmount: computed(() => {
        return actions.peg.calculateUnpegFee(Asset.get(symbol.value));
      }),
      handleBlur: () => {
        amount.value = trimZeros(amount.value);
      },
      handleSelectSymbol: () => {},
      handleMaxClicked: () => {
        if (!accountBalance.value) return;

        amount.value = accountBalance.value.toFixed();
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

      nextStepAllowed,
      nextStepMessage: computed(() => {
        return mode.value === "peg" ? "Peg" : "Unpeg";
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
        label="Amount"
      />
      <RaisedPanel>
        <RaisedPanelColumn v-if="mode === 'peg'">
          <Label>Sifchain Recipient Address</Label>
          <SifInput disabled v-model="address" />
        </RaisedPanelColumn>
        <RaisedPanelColumn v-if="mode === 'unpeg'">
          <Label>Ethereum Recipient Address</Label>
          <SifInput
            v-model="address"
            placeholder="Eg. 0xeaf65652e380528fffbb9fc276dd8ef608931e3c"
          />
        </RaisedPanelColumn>
      </RaisedPanel>
      <DetailsTable
        v-if="mode === 'unpeg'"
        :header="{
          show: amount !== '0.0',
          label: `${modeLabel} Amount`,
          data: `${amount} ${symbolLabel}`,
        }"
        :rows="[
          {
            show: !!feeAmount,
            label: 'Transaction Fee',
            data: `${feeAmount.toFixed(8)} cETH`,
          },
        ]"
      />
      <ActionsPanel
        connectType="connectToAll"
        @nextstepclick="handleActionClicked"
        :nextStepAllowed="nextStepAllowed"
        :nextStepMessage="nextStepMessage"
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
