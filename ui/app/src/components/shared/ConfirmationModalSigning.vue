<template>
  <div>
    <div class="confirmation">
      <div class="message">
        <Loader
          black
          :success="state === 'confirmed'"
          :failed="
            state === 'rejected' || state === 'failed' || state === 'out_of_gas'
          "
        /><br />
        <div class="text-wrapper">
          <!-- 
            TODO: This could be abstracted to AnimatedLoaderStateModal 
            that takes screens and switches them based on arbitrary state
            with arbitrary content that can be specified in page.
            The content below isn't really flexible enough and can be 
            templed into components

            Perhaps we could use render functions to accomplish this?
          -->
          <transition name="swipe">
            <div class="text" v-if="state === 'approving'">
              <p>Waiting for approval</p>
              <br />
              <p class="sub">Confirm this transaction in your wallet</p>
            </div>
          </transition>
          <transition name="swipe">
            <div class="text" v-if="state === 'signing'">
              <p>Waiting for confirmation</p>
              <slot name="signing"></slot>
              <br />
              <p class="sub">Confirm this transaction in your wallet</p>
            </div>
          </transition>

          <transition name="swipe">
            <div class="text" v-if="state === 'rejected'">
              <p>Transaction Rejected</p>
              <slot name="rejected"></slot>
              <br />
              <p class="sub">{{ transactionStateMsg }}</p>
            </div>
          </transition>

          <transition name="swipe">
            <div class="text" v-if="state === 'failed'">
              <p>Transaction Failed</p>
              <slot name="failed"></slot>
              <br />
              <p class="sub">{{ transactionStateMsg }}</p>
            </div>
          </transition>

          <transition name="swipe">
            <div class="text" v-if="state === 'out_of_gas'">
              <p>Transaction Failed - Out of Gas</p>
              <br />
              <p class="sub">Please try to increase the gas limit.</p>
            </div>
          </transition>

          <transition name="swipe">
            <div class="text" v-if="state === 'confirmed'">
              <p>Transaction Submitted</p>
              <slot name="confirmed"></slot>
              <br />
              <p class="sub">
                <!-- To the todo point above, we need to be able to control this better, hence isSifTxHash() -->
                <a
                  v-if="transactionHash?.substring(0, 2) !== '0x'"
                  class="anchor"
                  target="_blank"
                  :href="getBlockExplorerUrl(chainId, transactionHash)"
                  >View transaction on Block Explorer</a
                >
                <a
                  v-else
                  class="anchor"
                  target="_blank"
                  :href="`https://etherscan.io/tx/${transactionHash}`"
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
import { defineComponent } from "vue";
import { useCore } from "@/hooks/useCore";
import Loader from "@/components/shared/Loader.vue";
import SifButton from "@/components/shared/SifButton.vue";
import { getBlockExplorerUrl } from "./utils";

export default defineComponent({
  inheritAttrs: false,
  components: { Loader, SifButton },
  props: {
    state: String,
    transactionHash: String,
    transactionStateMsg: String,
  },

  setup() {
    const { config } = useCore();

    return {
      chainId: config.sifChainId,
      getBlockExplorerUrl,
    };
  },
});
</script>

<style lang="scss" scoped>
.confirmation {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
  padding: 15px 20px;
}
.message {
  width: 100%;
  font-size: 16px;
}
.text-wrapper {
  position: relative;
  display: flex;
  width: 100%;
  height: 88px;
}
.text {
  position: absolute;
  width: 100%;
}
.anchor {
  color: $c_black;
}
.sub {
  font-weight: normal;
  font-size: $fs_sm;
}
.footer {
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
