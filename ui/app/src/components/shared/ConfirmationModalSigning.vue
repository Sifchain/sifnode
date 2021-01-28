<template>
  <div>
    <div class="confirmation">
      <div class="message">
        <Loader black :success="confirmed" :failed="rejected || failed" /><br />
        <div class="text-wrapper">

          <transition name="swipe">
            <div class="text" v-if="signing">
              <p>Waiting for confirmation</p>
              <slot name="signing"></slot> 
              <br />
              <p class="sub">Confirm this transaction in your wallet</p>
            </div>
          </transition>
          
          <transition name="swipe">
            <div class="text" v-if="rejected">
              <p>Transaction Rejected</p>
              <slot name="rejected"></slot>
              <br />
              <p class="sub">{{ transactionStateMsg }}</p>
            </div>
          </transition>

          <transition name="swipe">
            <div class="text" v-if="failed">
              <p>Transaction Failed</p>
              <slot name="failed"></slot>
              <br />
              <p class="sub">{{ transactionStateMsg }}</p>
            </div>
          </transition>

          <transition name="swipe">
            <div class="text" v-if="confirmed">
              <p>Transaction Submitted</p>
              <slot name="confirmed"></slot>
              <br />
              <p class="sub">
                <a class="anchor" target="_blank" :href="`https://blockexplorer-${chainId}.sifchain.finance/transactions/${transactionHash}`"
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
import { computed } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import Loader from "@/components/shared/Loader.vue";
import SifButton from "@/components/shared/SifButton.vue";

export default defineComponent({
  components: { Loader, SifButton },
  props: {
    state: String,
    transactionHash: String,
    transactionStateMsg: String,
  },

  setup(props) {
    const { config } = useCore();

    const signing = computed(() => props.state === 'signing')
    const rejected = computed(() => props.state === 'rejected')
    const failed = computed(() => props.state === 'failed')
    const confirmed = computed(() => props.state === 'confirmed')

    return {
      chainId: config.sifChainId,
      signing,
      rejected,
      failed,
      confirmed,
    }
  }
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