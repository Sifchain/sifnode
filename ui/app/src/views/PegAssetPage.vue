<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { computed, ref, toRefs } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { Asset, AssetAmount } from "ui-core";
import CurrencyField from "@/components/currencyfield/CurrencyField.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";

import RaisedPanel from "@/components/shared/RaisedPanel.vue";
import { useRouter } from "vue-router";
import SifInput from "@/components/shared/SifInput.vue";
import DetailsTable from "@/components/shared/DetailsTable.vue";
import Label from "@/components/shared/Label.vue";
import RaisedPanelColumn from "@/components/shared/RaisedPanelColumn.vue";
import { trimZeros } from "ui-core/src/hooks/utils";
import BigNumber from "bignumber.js";
import {
  formatSymbol,
  getPeggedSymbol,
  getUnpeggedSymbol,
  useAssetItem,
} from "@/components/shared/utils";
import { toConfirmState } from "./utils/toConfirmState";
import { ConfirmState } from "../types";
import ConfirmationModal from "@/components/shared/ConfirmationModal.vue";
import { format } from "ui-core/src/utils/format";

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
    ConfirmationModal,
  },

  setup(props, context) {
    const { store, actions } = useCore();
    const router = useRouter();
    const mode = computed(() => {
      return router.currentRoute.value.path.indexOf("/peg/reverse") > -1
        ? "unpeg"
        : "peg";
    });

    const transactionState = ref<ConfirmState>("selecting");
    const transactionStateMsg = ref<string>("");
    const transactionHash = ref<string | null>(null);

    // const symbol = ref<string | null>(null);
    const symbol = computed(() => {
      const assetFrom = router.currentRoute.value.params.assetFrom;
      return Array.isArray(assetFrom) ? assetFrom[0] : assetFrom;
    });

    const oppositeSymbol = computed(() => {
      if (mode.value === "peg") {
        return getPeggedSymbol(symbol.value);
      }
      return getUnpeggedSymbol(symbol.value);
    });

    const amount = ref("0.0");
    const address = computed(() =>
      mode.value === "peg"
        ? store.wallet.sif.address
        : store.wallet.eth.address,
    );

    const isMaxActive = computed(() => {
      if (!accountBalance.value) return false;
      return (
        amount.value ===
        format(accountBalance.value.amount, accountBalance.value.asset)
      );
    });

    async function handlePegRequested() {
      const asset = Asset.get(symbol.value);
      if (asset.symbol !== "eth") {
        // if not eth you need to approve spend before peg
        transactionState.value = "approving";
        try {
          await actions.peg.approve(
            store.wallet.eth.address,
            AssetAmount(asset, amount.value),
          );
        } catch (err) {
          return (transactionState.value = "rejected");
        }
      }

      transactionState.value = "signing";
      const tx = await actions.peg.peg(AssetAmount(asset, amount.value));

      transactionHash.value = tx.hash;
      transactionState.value = toConfirmState(tx.state); // TODO: align states
      transactionStateMsg.value = tx.memo ?? "";
    }

    async function handleUnpegRequested() {
      transactionState.value = "signing";

      const tx = await actions.peg.unpeg(
        AssetAmount(Asset.get(symbol.value), amount.value),
      );

      transactionHash.value = tx.hash;
      transactionState.value = toConfirmState(tx.state); // TODO: align states
      transactionStateMsg.value = tx.memo ?? "";
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
      const balance =
        (accountBalance.value &&
          format(accountBalance.value.amount, accountBalance.value.asset)) ??
        "0.0";
      return (
        amountNum.isGreaterThan("0.0") &&
        address.value !== "" &&
        amountNum.isLessThanOrEqualTo(balance)
      );
    });

    function requestTransactionModalClose() {
      if (transactionState.value === "confirmed") {
        transactionState.value = "selecting";
        router.push("/peg"); // TODO push back to peg, but load unpeg tab when unpegging -> dynamic routing?
      } else {
        transactionState.value = "selecting";
      }
    }
    const feeAmount = computed(() => {
      return actions.peg.calculateUnpegFee(Asset.get(symbol.value));
    });

    const feeDisplayAmount = computed(() => {
      if (!feeAmount.value) return "";
      return format(feeAmount.value.amount, feeAmount.value.asset, {
        mantissa: 8,
      });
    });
    const pageState = {
      mode,
      modeLabel: computed(() => capitalize(mode.value)),
      symbol,
      symbolLabel: useAssetItem(symbol).label,
      amount,
      address,
      feeAmount,
      handleBlur: () => {
        if (isMaxActive.value === true) return;
        amount.value = trimZeros(amount.value);
      },
      handleSelectSymbol: () => {},
      handleMaxClicked: () => {
        if (!accountBalance.value) return;
        const decimals = Asset.get(symbol.value).decimals;
        const afterMaxValue =
          symbol.value === "ceth"
            ? accountBalance.value.subtract(feeAmount.value)
            : accountBalance.value;
        amount.value = afterMaxValue.lessThan("0")
          ? "0.0"
          : format(afterMaxValue, accountBalance.value.asset, {
              mantissa: decimals,
            });
      },
      handleAmountUpdated: (newAmount: string) => {
        amount.value = newAmount;
      },
      handleActionClicked: () => {
        transactionState.value = "confirming";
      },
      handlePegRequested,
      handleUnpegRequested,
      oppositeSymbol,
      formatSymbol,
      requestTransactionModalClose,
      transactionState,
      transactionStateMsg,
      transactionHash,
      nextStepAllowed,
      isMaxActive,
      feeDisplayAmount,
      nextStepMessage: computed(() => {
        return mode.value === "peg" ? "Peg" : "Unpeg";
      }),
    };
    (window as any).pageState = pageState;
    return pageState;
  },
});
</script>

<template>
  <Layout :title="mode === 'peg' ? 'Peg Asset' : 'Unpeg Asset'" backLink="/peg">
    <div class="vspace">
      <CurrencyField
        :amount="amount"
        :max="true"
        :isMaxActive="isMaxActive"
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
            disabled
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
            show: !!feeDisplayAmount,
            label: 'Transaction Fee',
            data: `${feeDisplayAmount} cETH`,
            tooltipMessage: `This is a fixed fee amount. This is a temporary solution as we are working towards improving this amount in upcoming versions of the network.`,
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
    <ConfirmationModal
      v-if="mode === 'peg'"
      @confirmed="handlePegRequested"
      :requestClose="requestTransactionModalClose"
      :state="transactionState"
      :transactionHash="transactionHash"
      :transactionStateMsg="transactionStateMsg"
      confirmButtonText="Confirm Peg"
      :title="`Peg token to Sifchain`"
    >
      <template v-slot:selecting>
        <DetailsTable
          :header="{
            show: amount !== '0.0',
            label: `${modeLabel} Amount`,
            data: `${amount} ${formatSymbol(symbol)}`,
          }"
          :rows="[
            {
              show: true,
              label: 'Direction',
              data: `${formatSymbol(symbol)} → ${formatSymbol(oppositeSymbol)}`,
            },
          ]"
        />
        <br />
        <p class="text--normal">
          *Please note your funds will be available for use on Sifchain only
          after 50 Ethereum block confirmations. This can take upwards of 20
          minutes.
        </p>
      </template>
      <template v-slot:approving>
        <p>Approving</p>
      </template>
      <template v-slot:common>
        <p class="text--normal">
          Pegging <span class="text--bold">{{ amount }} {{ symbol }}</span>
        </p>
      </template>
    </ConfirmationModal>
    <ConfirmationModal
      v-if="mode === 'unpeg'"
      @confirmed="handleUnpegRequested"
      :requestClose="requestTransactionModalClose"
      :state="transactionState"
      :transactionHash="transactionHash"
      :transactionStateMsg="transactionStateMsg"
      confirmButtonText="Confirm Unpeg"
      title="Unpeg token from Sifchain"
    >
      <template v-slot:selecting>
        <DetailsTable
          :header="{
            show: amount !== '0.0',
            label: `${modeLabel} Amount`,
            data: `${amount} ${formatSymbol(symbol)}`,
          }"
          :rows="[
            {
              show: true,
              label: 'Direction',
              data: `${formatSymbol(symbol)} → ${formatSymbol(oppositeSymbol)}`,
            },
            {
              show: !!feeDisplayAmount,
              label: 'Transaction Fee',
              data: `${feeDisplayAmount} cETH`,
            },
          ]"
        />
      </template>
      <template v-slot:common>
        <p class="text--normal">
          Unpegging <span class="text--bold">{{ amount }} {{ symbol }}</span>
        </p>
      </template>
    </ConfirmationModal>
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
