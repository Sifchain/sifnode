<template>
  <div>
    <div class="confirmation">
      <div class="message">
        <Loader
          black
          :success="state === 'success'"
          :failed="state === 'fail'"
        /><br />
        <div class="text-wrapper">
          <transition name="swipe">
            <div class="text" v-if="dialogState === 'submit'">
              <p>Waiting for confirmation</p>
              <p class="thin" data-handle="swap-message">
                Swapping
                <span class="thick">{{ _fromAmount }} {{ _fromToken }}</span>
                for
                <span class="thick">{{ _toAmount }} {{ _toToken }}</span>
              </p>
              <br />
              <p class="sub">Confirm this transaction in your wallet</p>
            </div>
          </transition>
          <transition name="swipe">
            <div class="text" v-if="dialogState === 'out_of_gas'">
              <p>Transaction Failed</p>
              <p class="thin" data-handle="swap-message">
                Failed to swap
                <span class="thick">{{ _fromAmount }} {{ _fromToken }}</span>
                for
                <span class="thick">{{ _toAmount }} {{ _toToken }}</span>
              </p>
              <br />
              <p class="sub">Please try to increase the gas limit.</p>
            </div>
          </transition>
          <transition name="swipe">
            <div class="text" v-if="dialogState === 'rejected'">
              <p>Transaction Rejected</p>
              <p class="thin" data-handle="swap-message">
                Failed to swap
                <span class="thick">{{ _fromAmount }} {{ _fromToken }}</span>
                for
                <span class="thick">{{ _toAmount }} {{ _toToken }}</span>
              </p>
              <br />
              <p class="sub">Please confirm the transaction in your wallet.</p>
            </div>
          </transition>
          <transition name="swipe">
            <div class="text" v-if="dialogState === 'slippage'">
              <p>Transaction Failed</p>
              <p class="thin" data-handle="swap-message">
                Failed to swap
                <span class="thick">{{ _fromAmount }} {{ _fromToken }}</span>
                for
                <span class="thick">{{ _toAmount }} {{ _toToken }}</span>
              </p>
              <br />
              <p class="sub">Please try to increase slippage tolerance.</p>
            </div>
          </transition>
          <transition name="swipe">
            <div class="text" v-if="dialogState === 'fail'">
              <p>Transaction Failed</p>
              <p class="thin" data-handle="swap-message">
                Failed to swap
                <span class="thick">{{ _fromAmount }} {{ _fromToken }}</span>
                for
                <span class="thick">{{ _toAmount }} {{ _toToken }}</span>
              </p>
              <br />
              <p class="sub"></p>
            </div>
          </transition>
          <transition name="swipe">
            <div class="text" v-if="dialogState === 'success'">
              <p>Transaction Submitted</p>
              <p class="thin" data-handle="swap-message">
                Swapped
                <span class="thick">{{ _fromAmount }} {{ _fromToken }}</span>
                for
                <span class="thick">{{ _toAmount }} {{ _toToken }}</span>
              </p>
              <br />
              <p class="sub">
                <a
                  class="anchor"
                  target="_blank"
                  :href="getBlockExplorerUrl(chainId, txStatus.hash)"
                  >View transaction on Block Explorer</a
                >
              </p>
            </div>
          </transition>
        </div>
      </div>
    </div>
    <div class="footer" :class="{ confirmed }">
      <SifButton block @click="$emit('closerequested')" primary
        >Close</SifButton
      >
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from "vue";
import { computed } from "@vue/reactivity";
import Loader from "@/components/shared/Loader.vue";
import SifButton from "@/components/shared/SifButton.vue";
import { useCore } from "@/hooks/useCore";
import { getBlockExplorerUrl } from "../shared/utils";
import { ErrorCode, TransactionStatus } from "ui-core";
import { UiState } from "../../views/SwapPage.vue";

export default defineComponent({
  components: { Loader, SifButton },
  props: {
    txStatus: { type: Object as PropType<TransactionStatus>, default: null },
    confirmed: Boolean,
    failed: Boolean,
    state: { type: String as PropType<UiState> },
    fromAmount: String,
    fromToken: String,
    toAmount: String,
    toToken: String,
  },
  setup(props) {
    const { config } = useCore();

    // Unfortunately because of the way vue works we need to
    // tokenize all these states here instead of writing them inline
    // This is annoying as it means we have more chance of error
    // Best way to fix this is to explore strategies for converting to JSX
    const dialogState = computed(() => {
      if (props.state === "submit") return "submit";
      if (props.state === "success") return "success";
      if (
        props.state === "fail" &&
        props.txStatus.code === ErrorCode.TX_FAILED_OUT_OF_GAS
      )
        return "out_of_gas";

      if (
        props.state === "fail" &&
        props.txStatus.code === ErrorCode.USER_REJECTED
      )
        return "rejected";

      if (
        props.state === "fail" &&
        props.txStatus.code === ErrorCode.TX_FAILED_SLIPPAGE
      )
        return "slippage";

      if (props.state === "fail") return "unknown";
    });

    // Need to cache amounts and disconnect reactivity
    return {
      _fromAmount: props.fromAmount,
      _fromToken: props.fromToken,
      _toAmount: props.toAmount,
      _toToken: props.toToken,
      chainId: config.sifChainId,
      getBlockExplorerUrl,
      ErrorCode,
      dialogState,
    };
  },
});
</script>

<style lang="scss" scoped>
.confirmation {
  display: flex;
  justify-content: start;
  align-items: start;
  min-height: 40vh;
  padding: 15px 20px;
}
.message {
  width: 100%;
  font-size: 18px;
  margin-top: 3em;
}
.text-wrapper {
  margin-top: 0.5em;
  position: relative;
  display: flex;
  width: 100%;
}
.text {
  position: absolute;
  width: 100%;
}
.anchor {
  color: $c_black;
}
.thin {
  margin-top: 1em;
  margin-bottom: 1em;
  font-size: 15px;
  font-weight: normal;
}
.thick {
  font-weight: bold;
}
.sub {
  font-weight: normal;
  font-size: $fs_sm;
}
.footer {
  padding: 16px;
  visibility: hidden;
  transition: opacity 0.5s ease-out;
  opacity: 0;
  &.confirmed {
    opacity: 1;
    visibility: inherit;
  }
}
.swipe-enter-active,
.swipe-leave-active {
  transition: transform 0.5s ease-out;
}

.swipe-enter-from {
  transform: translateX(100%);
}
.swipe-leave-to {
  transform: translateX(-100%);
}
</style>
