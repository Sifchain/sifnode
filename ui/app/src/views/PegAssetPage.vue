<script lang="tsx">
import { defineComponent, Slot, watch } from "vue";
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
import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationModalAsk from "../components/shared/ConfirmationModalAsk.vue";
import ConfirmationModalSwipe from "../components/shared/ConfirmationModalSwipe.vue";
import SwipeMessage from "../components/shared/ConfirmationModalSwipeMessage.vue";
import VSpace from "../components/shared/VSpace.vue";

// Feel like it is overkill now but it would be worth looking
// into some kind of statemachine to manage these flows?
// Xstate gives you some nice debuggig tools but is verbose and doesn't work well with TS
// Alternatively we could make things router based and use that as a statemachine
type PageStates =
  | "idle"
  | "confirm"
  | "approve"
  | "sign"
  | "success"
  | "fail"
  | "reject";

// This is a little TS utility function for validating our
// slots match the states we allow
function validSlots<T extends string>(o: { [k in T]?: any }) {
  return o;
}

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export default defineComponent({
  setup() {
    const { store, actions, config } = useCore();
    const router = useRouter();

    const pageState = ref<PageStates>("idle");

    const mode = computed(() => {
      return router.currentRoute.value.path.indexOf("/peg/reverse") > -1
        ? "unpeg"
        : "peg";
    });

    watch(pageState, (newState, prevState) => {
      // When we are moving from success to idle head back to peg
      if (prevState === "success" && newState === "idle") {
        router.push("/peg");
      }
    });

    const transactionHash = ref<string | null>(null);
    const transactionErrorMessage = ref<string | null>(null);

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
      pageState.value = "approve";
      const assetAmount = AssetAmount(Asset.get(symbol.value), amount.value);

      try {
        await actions.peg.approveSpend(assetAmount);
        pageState.value = "sign";
      } catch (err) {
        pageState.value = "reject";
        transactionErrorMessage.value =
          "You failed to approve funds management.";
        return;
      }

      const tx = await actions.peg.peg(assetAmount);

      if (!tx.hash) {
        pageState.value = "fail";
        return;
      }

      transactionHash.value = tx.hash;

      pageState.value = "success";
    }

    async function handleUnpegRequested() {
      pageState.value = "sign";

      const tx = await actions.peg.unpeg(
        AssetAmount(Asset.get(symbol.value), amount.value)
      );

      if (!tx.hash) {
        pageState.value = "fail";
        return;
      }

      transactionHash.value = tx.hash;
      pageState.value = "success";
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
      pageState.value = "confirm";
    }

    function handleCloseClicked() {
      pageState.value = "idle";
    }

    function handleActionConfirmed() {
      if (mode.value === "peg") {
        handlePegRequested();
      } else {
        handleUnpegRequested();
      }
    }

    const nextStepMessage = computed(() => {
      return mode.value === "peg" ? "Peg" : "Unpeg";
    });

    return () => {
      const peggingMessage =
        mode.value === "peg" ? (
          <p class="text--normal">
            Pegging{" "}
            <span class="text--bold">
              {amount.value} {symbol.value}
            </span>
          </p>
        ) : (
          <p class="text--normal">
            Unpegging{" "}
            <span class="text--bold">
              {amount.value} {symbol.value}
            </span>
          </p>
        );

      return (
        <Layout
          title={mode.value === "peg" ? "Peg Asset" : "Unpeg Asset"}
          backLink="/peg"
        >
          <VSpace>
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
                  <SifInput disabled v-model={address.value} />
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
          </VSpace>
          <ModalView
            requestClose={handleCloseClicked}
            isOpen={pageState.value !== "idle"}
          >
            {pageState.value === "confirm" ? (
              <ConfirmationModalAsk
                confirmButtonText="Confirm Peg"
                onConfirmed={handleActionConfirmed}
                title={
                  mode.value === "peg"
                    ? "Peg token to Sifchain"
                    : "Unpeg token from Sifchain"
                }
              >
                <>
                  <DetailsTable
                    header={{
                      show: amount.value !== "0.0",
                      label: `${modeLabel.value} Amount`,
                      data: `${amount.value} ${formatSymbol(symbol.value)}`,
                    }}
                    rows={
                      mode.value === "peg"
                        ? [
                            {
                              show: true,
                              label: "Direction",
                              data: `${formatSymbol(
                                symbol.value
                              )} → ${formatSymbol(oppositeSymbol.value)}`,
                            },
                          ]
                        : [
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
                          ]
                    }
                  />
                  {mode.value === "peg" && (
                    <>
                      <br />
                      <p class="text--normal">
                        *Please note your funds will be available for use on
                        Sifchain only after 50 Ethereum block confirmations.
                        This can take upwards of 20 minutes.
                      </p>
                    </>
                  )}
                </>
              </ConfirmationModalAsk>
            ) : (
              <ConfirmationModalSwipe
                state={pageState.value}
                loaderState={{
                  success: { success: true, failed: false },
                  fail: { success: false, failed: true },
                  reject: { success: false, failed: true },
                }}
                v-slots={validSlots<PageStates>({
                  approve: () => (
                    <SwipeMessage
                      title="Waiting for approval"
                      sub="Confirm this transaction in your wallet"
                    >
                      {peggingMessage}
                    </SwipeMessage>
                  ),
                  sign: () => (
                    <SwipeMessage
                      title="Waiting for confirmation"
                      sub="Confirm this transaction in your wallet"
                    >
                      {peggingMessage}
                    </SwipeMessage>
                  ),
                  success: () => (
                    <SwipeMessage
                      title="Transaction Submitted"
                      v-slots={{
                        sub: () =>
                          mode.value === "peg" ? (
                            <a
                              class="anchor"
                              target="_blank"
                              href={`https://blockexplorer-${config.sifChainId}.sifchain.finance/transactions/${transactionHash.value}`}
                            >
                              View transaction on Block Explorer
                            </a>
                          ) : (
                            <a
                              class="anchor"
                              target="_blank"
                              href={`https://etherscan.io/tx/${transactionHash.value}`}
                            >
                              View transaction on Block Explorer
                            </a>
                          ),
                      }}
                    >
                      {peggingMessage}
                    </SwipeMessage>
                  ),
                  fail: () => (
                    <SwipeMessage
                      title="Transaction Failed"
                      sub={transactionErrorMessage.value}
                    >
                      {peggingMessage}
                    </SwipeMessage>
                  ),
                  reject: () => (
                    <SwipeMessage
                      title="Transaction Rejected"
                      sub={transactionErrorMessage.value}
                    >
                      {peggingMessage}
                    </SwipeMessage>
                  ),
                })}
              />
            )}
          </ModalView>
        </Layout>
      );
    };
  },
});
</script>

