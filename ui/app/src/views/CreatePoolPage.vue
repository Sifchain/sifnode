<script lang="ts">
import { defineComponent, ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import Layout from "@/components/layout/Layout.vue";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import SelectTokenDialogSif from "@/components/tokenSelector/SelectTokenDialogSif.vue";
import Modal from "@/components/shared/Modal.vue";
import { PoolState, usePoolCalculator } from "ui-core";
import { useCore } from "@/hooks/useCore";
import { useWallet } from "@/hooks/useWallet";
import { computed } from "@vue/reactivity";
import FatInfoTable from "@/components/shared/FatInfoTable.vue";
import FatInfoTableCell from "@/components/shared/FatInfoTableCell.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import { useCurrencyFieldState } from "@/hooks/useCurrencyFieldState";
import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationModalAsk from "../components/shared/ConfirmationModalAsk.vue";
import ConfirmationModalSwipe from "../components/shared/ConfirmationModalSwipe.vue";
import SwipeMessage from "@/components/shared/ConfirmationModalSwipeMessage.vue";
import DetailsPanelPool from "@/components/shared/DetailsPanelPool.vue";
import { formatNumber } from "@/components/shared/utils";

type PageStates = "idle" | "confirm" | "sign" | "success" | "fail" | "reject";

export default defineComponent({
  components: {
    ActionsPanel,
    Layout,
    Modal,
    CurrencyPairPanel,
    SelectTokenDialogSif,
    ModalView,
    ConfirmationModalAsk,
    ConfirmationModalSwipe,
    SwipeMessage,
    DetailsPanelPool,
    FatInfoTable,
    FatInfoTableCell,
  },
  props: ["title"],
  setup() {
    const { actions, poolFinder, store } = useCore();
    const selectedField = ref<"from" | "to" | null>(null);
    const pageState = ref<PageStates>("idle");
    const errorMessage = ref<string | null>(null);
    const transactionHash = ref<string | null>(null);
    const router = useRouter();
    const route = useRoute();

    const { fromSymbol, fromAmount, toAmount } = useCurrencyFieldState();

    const toSymbol = ref("rowan");

    watch(pageState, (newState, prevState) => {
      // When we are moving from success to idle head back to pools
      if (prevState === "success" && newState === "idle") {
        router.push("/pool");
        clearAmounts();
      }
    });

    fromSymbol.value = route.params.externalAsset
      ? route.params.externalAsset.toString()
      : null;

    function clearAmounts() {
      fromAmount.value = "0.0";
      toAmount.value = "0.0";
    }

    const { connected, connectedText } = useWalletButton({
      addrLen: 8,
    });

    const { balances } = useWallet(store);

    const liquidityProvider = computed(() => {
      if (!fromSymbol) return null;
      return (
        store.accountpools.find((pool) => {
          return pool.lp.asset.symbol === fromSymbol.value;
        })?.lp ?? null
      );
    });

    const {
      aPerBRatioMessage,
      bPerARatioMessage,
      aPerBRatioProjectedMessage,
      bPerARatioProjectedMessage,
      shareOfPoolPercent,
      totalLiquidityProviderUnits,
      tokenAFieldAmount,
      tokenBFieldAmount,
      preExistingPool,
      state,
    } = usePoolCalculator({
      balances,
      tokenAAmount: fromAmount,
      tokenBAmount: toAmount,
      tokenASymbol: fromSymbol,
      tokenBSymbol: toSymbol,
      poolFinder,
      liquidityProvider,
    });

    function handleNextStepClicked() {
      if (!tokenAFieldAmount.value)
        throw new Error("from field amount is not defined");
      if (!tokenBFieldAmount.value)
        throw new Error("to field amount is not defined");

      pageState.value = "confirm";
    }

    async function handleAddLiquidityClicked() {
      if (!tokenAFieldAmount.value)
        throw new Error("Token A field amount is not defined");
      if (!tokenBFieldAmount.value)
        throw new Error("Token B field amount is not defined");

      pageState.value = "sign";
      const tx = await actions.clp.addLiquidity(
        tokenBFieldAmount.value,
        tokenAFieldAmount.value
      );

      transactionHash.value = tx.hash;

      if (tx.state === "failed") {
        pageState.value = "fail";
        errorMessage.value = tx.memo ?? "The transaction failed";
      }

      if (tx.state === "rejected") {
        pageState.value = "reject";
        errorMessage.value = tx.memo ?? "You rejected the transaction";
      }

      pageState.value = "success";
      errorMessage.value = "";
    }

    function requestTransactionModalClose() {
      pageState.value = "idle";
    }

    return {
      fromAmount,
      fromSymbol,

      toAmount,
      toSymbol,

      connected,
      aPerBRatioMessage,
      bPerARatioMessage,
      aPerBRatioProjectedMessage,
      bPerARatioProjectedMessage,
      nextStepMessage: computed(() => {
        switch (state.value) {
          case PoolState.SELECT_TOKENS:
            return "Select Tokens";
          case PoolState.ZERO_AMOUNTS:
            return "Please enter an amount";
          case PoolState.INSUFFICIENT_FUNDS:
            return "Insufficient Funds";
          case PoolState.VALID_INPUT:
            return preExistingPool.value ? "Add liquidity" : "Create Pool";
        }
      }),
      nextStepAllowed: computed(() => {
        return state.value === PoolState.VALID_INPUT;
      }),
      handleFromSymbolClicked(next: () => void) {
        selectedField.value = "from";
        next();
      },
      handleToSymbolClicked(next: () => void) {
        selectedField.value = "to";
        next();
      },
      handleSelectClosed(data: string) {
        if (typeof data !== "string") {
          return;
        }

        if (selectedField.value === "from") {
          fromSymbol.value = data;
        }

        if (selectedField.value === "to") {
          toSymbol.value = data;
        }
        selectedField.value = null;
      },

      handleNextStepClicked,

      handleAddLiquidityClicked,

      transactionHash,

      requestTransactionModalClose,

      pageState,
      errorMessage,
      handleBlur() {
        selectedField.value = null;
      },
      handleFromFocused() {
        selectedField.value = "from";
      },
      handleToFocused() {
        selectedField.value = "to";
      },
      handleFromMaxClicked() {
        selectedField.value = "from";
        const accountBalance = balances.value.find(
          (balance) => balance.asset.symbol === fromSymbol.value
        );
        if (!accountBalance) return;
        fromAmount.value = accountBalance.toFixed(8);
      },
      handleToMaxClicked() {
        selectedField.value = "to";
        const accountBalance = balances.value.find(
          (balance) => balance.asset.symbol === toSymbol.value
        );
        if (!accountBalance) return;
        toAmount.value = accountBalance.toFixed(8);
      },
      shareOfPoolPercent,
      connectedText,
      formatNumber,

      poolUnits: totalLiquidityProviderUnits,
    };
  },
});
</script>

<template>
  <Layout
    class="pool"
    :backLink="`${fromSymbol ? '/pool/' + fromSymbol : '/pool'}`"
    :title="title"
  >
    <Modal @close="handleSelectClosed">
      <template v-slot:activator="{ requestOpen }">
        <CurrencyPairPanel
          v-model:fromAmount="fromAmount"
          v-model:fromSymbol="fromSymbol"
          @fromfocus="handleFromFocused"
          @fromblur="handleBlur"
          @fromsymbolclicked="handleFromSymbolClicked(requestOpen)"
          :fromSymbolSelectable="connected"
          :fromMax="true"
          @frommaxclicked="handleFromMaxClicked"
          v-model:toAmount="toAmount"
          v-model:toSymbol="toSymbol"
          @tofocus="handleToFocused"
          @toblur="handleBlur"
          :toMax="true"
          @tomaxclicked="handleToMaxClicked"
          toSymbolFixed
          canSwapIcon="plus"
      /></template>
      <template v-slot:default="{ requestClose }">
        <SelectTokenDialogSif
          :selectedTokens="[fromSymbol, toSymbol].filter(Boolean)"
          @tokenselected="requestClose"
        />
      </template>
    </Modal>

    <FatInfoTable :show="nextStepAllowed">
      <template #header>Pool Token Prices</template>
      <template #body>
        <FatInfoTableCell>
          <span class="number">{{ formatNumber(aPerBRatioMessage) }}</span
          ><br />
          <span
            >{{ fromSymbol.toUpperCase() }} per
            {{ toSymbol.toUpperCase() }}</span
          >
        </FatInfoTableCell>
        <FatInfoTableCell>
          <span class="number">{{ formatNumber(bPerARatioMessage) }}</span
          ><br />
          <span
            >{{ toSymbol.toUpperCase() }} per
            {{ fromSymbol.toUpperCase() }}</span
          > </FatInfoTableCell
        ><FatInfoTableCell />
      </template>
    </FatInfoTable>

    <FatInfoTable :show="nextStepAllowed">
      <template #header>Price Impact and Pool Share</template>
      <template #body>
        <FatInfoTableCell>
          <span class="number">{{
            formatNumber(aPerBRatioProjectedMessage)
          }}</span
          ><br />
          <span
            >{{ fromSymbol.toUpperCase() }} per
            {{ toSymbol.toUpperCase() }}</span
          >
        </FatInfoTableCell>
        <FatInfoTableCell>
          <span class="number">{{
            formatNumber(bPerARatioProjectedMessage)
          }}</span
          ><br />
          <span
            >{{ toSymbol.toUpperCase() }} per
            {{ fromSymbol.toUpperCase() }}</span
          >
        </FatInfoTableCell>
        <FatInfoTableCell>
          <span class="number">{{ shareOfPoolPercent }}</span
          ><br />Share of Pool
        </FatInfoTableCell></template
      >
    </FatInfoTable>

    <ActionsPanel
      @nextstepclick="handleNextStepClicked"
      :nextStepAllowed="nextStepAllowed"
      :nextStepMessage="nextStepMessage"
    />
    <ModalView
      :isOpen="pageState !== 'idle'"
      :requestClose="requestTransactionModalClose"
    >
      <ConfirmationModalAsk
        v-if="pageState === 'confirm'"
        confirmButtonText="Create Pool"
        :onConfirmed="handleAddLiquidityClicked"
        :title="
          mode === 'peg' ? 'Peg token to Sifchain' : 'Unpeg token from Sifchain'
        "
        ><div>
          <DetailsPanelPool
            class="details"
            :fromTokenLabel="fromSymbol"
            :fromAmount="fromAmount"
            :toTokenLabel="toSymbol"
            :toAmount="toAmount"
            :aPerB="aPerBRatioMessage"
            :bPerA="bPerARatioMessage"
            :shareOfPool="shareOfPoolPercent"
          />
        </div>
      </ConfirmationModalAsk>
      <ConfirmationModalSwipe
        v-else
        :state="pageState"
        :loaderState="{
          success: { success: true, failed: false },
          fail: { success: false, failed: true },
          reject: { success: false, failed: true },
        }"
      >
        <template #sign>
          <SwipeMessage
            title="Waiting for confirmation"
            sub="Confirm this transaction in your wallet"
          >
            <p class="text--normal">
              Supplying
              <span class="text--bold">{{ fromAmount }} {{ fromSymbol }}</span>
              and
              <span class="text--bold">{{ toAmount }} {{ toSymbol }}</span>
            </p>
          </SwipeMessage>
        </template>
        <template #success>
          <SwipeMessage title="Transaction Submitted"
            ><template #sub>
              <a
                v-if="mode === 'peg'"
                class="anchor"
                target="_blank"
                :href="`https://blockexplorer-${config.sifChainId}.sifchain.finance/transactions/${transactionHash}`"
              >
                View transaction on Block Explorer
              </a>

              <a
                v-else
                class="anchor"
                target="_blank"
                :href="`https://etherscan.io/tx/${transactionHash}`"
              >
                View transaction on Block Explorer
              </a>
            </template>

            <template #default
              ><p class="text--normal">
                Supplying
                <span class="text--bold"
                  >{{ fromAmount }} {{ fromSymbol }}</span
                >
                and
                <span class="text--bold">{{ toAmount }} {{ toSymbol }}</span>
              </p></template
            >
          </SwipeMessage>
        </template>
        <template #fail>
          <SwipeMessage title="Transaction Failed" :sub="errorMessage">
            <p class="text--normal">
              Supplying
              <span class="text--bold">{{ fromAmount }} {{ fromSymbol }}</span>
              and
              <span class="text--bold">{{ toAmount }} {{ toSymbol }}</span>
            </p>
          </SwipeMessage>
        </template>
        <template #reject>
          <SwipeMessage title="Transaction Rejected" :sub="errorMessage">
            <p class="text--normal">
              Supplying
              <span class="text--bold">{{ fromAmount }} {{ fromSymbol }}</span>
              and
              <span class="text--bold">{{ toAmount }} {{ toSymbol }}</span>
            </p>
          </SwipeMessage>
        </template>
      </ConfirmationModalSwipe>
    </ModalView>
  </Layout>
</template>

<style lang="scss" scoped>
.number {
  font-size: 16px;
  font-weight: bold;
}
.anchor {
  color: $c_black;
}
</style>
