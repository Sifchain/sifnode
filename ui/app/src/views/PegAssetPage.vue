<script lang="tsx">
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

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export default defineComponent({
  setup(props, context) {
    const { store, actions } = useCore();
    const router = useRouter();
    const mode = computed(() => {
      return router.currentRoute.value.path.indexOf("/peg/reverse") > -1
        ? "unpeg"
        : "peg";
    });

    const transactionState = ref<ConfirmState>("selecting");
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
      mode.value === "peg" ? store.wallet.sif.address : store.wallet.eth.address
    );

    async function handlePegRequested() {
      // transactionState.value = "signing";
      const tx = await actions.peg.peg(
        AssetAmount(Asset.get(symbol.value), amount.value)
      );

      transactionHash.value = tx.hash;
      // transactionState.value = toConfirmState(tx.state); // TODO: align states
      // transactionStateMsg.value = tx.memo ?? "";
    }

    async function handleUnpegRequested() {
      // transactionState.value = "signing";

      const tx = await actions.peg.unpeg(
        AssetAmount(Asset.get(symbol.value), amount.value)
      );

      transactionHash.value = tx.hash;
      // transactionState.value = toConfirmState(tx.state); // TODO: align states
      // transactionStateMsg.value = tx.memo ?? "";
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
      const balance = accountBalance.value?.toFixed(18) ?? "0.0";
      return (
        amountNum.isGreaterThan("0.0") &&
        address.value !== "" &&
        amountNum.isLessThanOrEqualTo(balance)
      );
    });

    const currentTxStatus = computed(() => {
      if (!transactionHash.value) return null;
      return store.tx.hash[transactionHash.value] ?? null;
    });

    function requestTransactionModalClose() {
      if (
        currentTxStatus.value &&
        toConfirmState(currentTxStatus.value?.state) === "confirmed"
      ) {
        confirmationModalOpen.value = false;
        router.push("/peg"); // TODO push back to peg, but load unpeg tab when unpegging -> dynamic routing?
      } else {
        confirmationModalOpen.value = false;
      }
    }

    const confirmationModalOpen = ref(false);
    function handleMaxClicked() {
      if (!accountBalance.value) return;

      amount.value = accountBalance.value.toFixed();
    }

    function handleBlur() {
      amount.value = trimZeros(amount.value);
    }

    function handleAmountUpdated(newAmount: string) {
      amount.value = newAmount;
    }
    const modeLabel = computed(() => capitalize(mode.value));
    const symbolLabel = useAssetItem(symbol).label;

    const feeAmount = computed(() => {
      return actions.peg.calculateUnpegFee(Asset.get(symbol.value));
    });

    function handleActionClicked() {
      transactionState.value = "confirming";
      confirmationModalOpen.value = true;
    }

    const nextStepMessage = computed(() => {
      return mode.value === "peg" ? "Peg" : "Unpeg";
    });

    return () => {
      return (
        <Layout
          title={mode.value === "peg" ? "Peg Asset" : "Unpeg Asset"}
          backLink="/peg"
        >
          <div class="vspace">
            <CurrencyField
              amount={amount.value}
              {...{ "onUpdate:amount": handleAmountUpdated }}
              max
              selectable
              symbol={symbol.value}
              symbolFixed
              onBlur={handleBlur}
              onMaxclicked={handleMaxClicked}
              label="Amount"
            />
            <RaisedPanel>
              {mode.value === "peg" ? (
                <RaisedPanelColumn>
                  <Label>Sifchain Recipient Address</Label>
                  <SifInput disabled v-model={address} />
                </RaisedPanelColumn>
              ) : (
                <RaisedPanelColumn>
                  <Label>Ethereum Recipient Address</Label>
                  <SifInput
                    v-model={address.value}
                    placeholder="Eg. 0xeaf65652e380528fffbb9fc276dd8ef608931e3c"
                  />
                </RaisedPanelColumn>
              )}
            </RaisedPanel>
            {mode.value === "unpeg" && (
              <DetailsTable
                header={{
                  show: amount.value !== "0.0",
                  label: `${modeLabel.value} Amount`,
                  data: `${amount.value} ${symbolLabel.value}`,
                }}
                rows={[
                  {
                    show: !!feeAmount.value,
                    label: "Transaction Fee",
                    data: `${feeAmount.value.toFixed(8)} cETH`,
                  },
                ]}
              />
            )}
            <ActionsPanel
              connectType="connectToAll"
              onNextstepclick={handleActionClicked}
              nextStepAllowed={nextStepAllowed.value}
              nextStepMessage={nextStepMessage.value}
            />
            {mode.value === "peg" && (
              <ConfirmationModal
                isOpen={confirmationModalOpen.value}
                onConfirmed={handlePegRequested}
                requestClose={requestTransactionModalClose}
                txStatus={currentTxStatus.value ?? undefined}
                confirmButtonText="Confirm Peg"
                title="Peg token to Sifchain"
                v-slots={{
                  common: () => (
                    <p class="text--normal">
                      Pegging{" "}
                      <span class="text--bold">
                        {{ amount }} {{ symbol }}
                      </span>
                    </p>
                  ),
                  selecting: () => (
                    <>
                      <DetailsTable
                        header={{
                          show: amount.value !== "0.0",
                          label: `${modeLabel.value} Amount`,
                          data: `${amount.value} ${formatSymbol(symbol.value)}`,
                        }}
                        rows={[
                          {
                            show: true,
                            label: "Direction",
                            data: `${formatSymbol(
                              symbol.value
                            )} → ${formatSymbol(oppositeSymbol.value)}`,
                          },
                        ]}
                      />
                      <br />
                      <p class="text--normal">
                        *Please note your funds will be available for use on
                        Sifchain only after 50 Ethereum block confirmations.
                        This can take upwards of 20 minutes.
                      </p>
                    </>
                  ),
                }}
              ></ConfirmationModal>
            )}

            {mode.value === "unpeg" && (
              <ConfirmationModal
                isOpen={confirmationModalOpen.value}
                onConfirmed={handlePegRequested}
                requestClose={requestTransactionModalClose}
                txStatus={currentTxStatus.value ?? undefined}
                confirmButtonText="Confirm Unpeg"
                title="Unpeg token from Sifchain"
                v-slots={{
                  common: () => (
                    <p class="text--normal">
                      Unpegging{" "}
                      <span class="text--bold">
                        {{ amount }} {{ symbol }}
                      </span>
                    </p>
                  ),
                  selecting: () => (
                    <>
                      <DetailsTable
                        header={{
                          show: amount.value !== "0.0",
                          label: `${modeLabel.value} Amount`,
                          data: `${amount.value} ${formatSymbol(symbol.value)}`,
                        }}
                        rows={[
                          {
                            show: true,
                            label: "Direction",
                            data: `${formatSymbol(
                              symbol.value
                            )} → ${formatSymbol(oppositeSymbol.value)}`,
                          },
                          {
                            show: !!feeAmount,
                            label: "Transaction Fee",
                            data: `${feeAmount.value.toFixed(8)} cETH`,
                          },
                        ]}
                      />
                    </>
                  ),
                }}
              ></ConfirmationModal>
            )}
          </div>
        </Layout>
      );
    };
  },
});
</script>

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
