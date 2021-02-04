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
import { createMachine, interpret } from "xstate";
import ConfirmationModalAsk from "../components/shared/ConfirmationModalAsk.vue";
import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationModalSwipe from "../components/shared/ConfirmationModalSwipe.vue";
import SwipeMessage from "../components/shared/ConfirmationModalSwipeMessage.vue";

type Context = any;

type StateEvent<T> = { type: T };

type StateEvents =
  | StateEvent<"ACTION_CLICKED">
  | StateEvent<"UNPEG_REQUESTED">
  | StateEvent<"PEG_REQUESTED">
  | StateEvent<"USER_APPROVED_SPEND">
  | StateEvent<"USER_REJECTED">
  | StateEvent<"SUBMITTED">
  | StateEvent<"FAIL">
  | StateEvent<"SUCCESS">
  | StateEvent<"CLOSE_CLICKED">;

// This state machine defines the sequence and possible states of this page
// and allows for a clearer way to think about the possible state of the component
const pegAssetsStateMachine = createMachine<Context, StateEvents>({
  initial: "idle",
  states: {
    idle: { on: { ACTION_CLICKED: "confirm" } },
    confirm: {
      on: {
        UNPEG_REQUESTED: "sign",
        PEG_REQUESTED: "approve",
        CLOSE_CLICKED: "idle",
      },
    },
    approve: {
      on: {
        USER_APPROVED_SPEND: "sign",
        CLOSE_CLICKED: "idle",
        USER_REJECTED: "reject",
      },
    },
    sign: {
      on: {
        SUBMITTED: "submit",
        CLOSE_CLICKED: "idle",
        USER_REJECTED: "reject",
      },
    },
    submit: {
      on: {
        FAIL: "fail",
        SUCCESS: "success",
        CLOSE_CLICKED: "idle",
      },
    },
    success: { on: { CLOSE_CLICKED: "idle" } },
    fail: { on: { CLOSE_CLICKED: "idle" } },
    reject: { on: { CLOSE_CLICKED: "idle" } },
  },
});

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export default defineComponent({
  setup(props, context) {
    const { store, actions, config } = useCore();
    const router = useRouter();
    const pageController = interpret(pegAssetsStateMachine, {
      devTools: true,
    });

    pageController.start();
    const pageState = ref(pageController.state.value.toString());

    pageController.onTransition((s, x) => {
      const stateName = s.value.toString();
      const prevPageState = pageState.value;
      pageState.value = stateName;
      if (prevPageState === "success" && stateName === "idle") {
        router.push("/peg");
      }
    });

    const mode = computed(() => {
      return router.currentRoute.value.path.indexOf("/peg/reverse") > -1
        ? "unpeg"
        : "peg";
    });

    const transactionHash = ref<string | null>(null);
    const transactionMessage = ref<string | undefined>(undefined);
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
      pageController.send("PEG_REQUESTED");
      const assetAmount = AssetAmount(Asset.get(symbol.value), amount.value);

      try {
        await actions.peg.approveSpend(assetAmount);
        pageController.send("USER_APPROVED_SPEND");
      } catch (err) {
        pageController.send("USER_REJECTED");
        return;
      }

      const tx = await actions.peg.peg(assetAmount);
      pageController.send("SUBMITTED");

      if (!tx.hash) {
        pageController.send("FAIL");
        return;
      }

      transactionHash.value = tx.hash;
      transactionMessage.value = tx.memo;
      pageController.send("SUCCESS");
    }

    async function handleUnpegRequested() {
      pageController.send("UNPEG_REQUESTED");

      const tx = await actions.peg.unpeg(
        AssetAmount(Asset.get(symbol.value), amount.value)
      );
      pageController.send("SUBMITTED");

      if (!tx.hash) {
        pageController.send("FAIL");
        return;
      }

      transactionHash.value = tx.hash;
      pageController.send("SUCCESS");
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
      pageController.send("ACTION_CLICKED");
    }

    function handleCloseClicked() {
      pageController.send("CLOSE_CLICKED");
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
          </div>
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
                v-slots={{
                  body: () => (
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
                  ),
                }}
              />
            ) : (
              <ConfirmationModalSwipe
                state={pageState.value}
                loaderState={{
                  success: { success: true, failed: false },
                  fail: { success: false, failed: true },
                  reject: { success: false, failed: true },
                }}
                v-slots={{
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
                      sub={transactionMessage.value}
                    >
                      {peggingMessage}
                    </SwipeMessage>
                  ),
                  reject: () => (
                    <SwipeMessage
                      title="Transaction Rejected"
                      sub={transactionMessage.value}
                    >
                      {peggingMessage}
                    </SwipeMessage>
                  ),
                }}
              />
            )}
          </ModalView>
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
// XXX: Add to ConfirmModal
.modal-inner {
  display: flex;
  flex-direction: column;
  padding: 30px 20px 20px 20px;
  min-height: 50vh;
}
</style>
