<template>
  <div>
    <div class="confirmation">
      <div class="message">
        <Loader black :success="confirmed" /><br />
        <div class="text-wrapper">
          <transition name="swipe">
            <div class="text" v-if="state === 'signing'">
              <p>Waiting for confirmation</p>
              <p class="thin">
                Supplying
                <span class="thick">{{ _fromAmount }} {{ _fromToken }}</span>
                and
                <span class="thick">{{ _toAmount }} {{ _toToken }}</span>
              </p>
              <br />
              <p class="sub">Confirm this transaction in your wallet</p>
            </div>
          </transition>
          <transition name="swipe">
            <div class="text" v-if="confirmed">
              <p>Transaction Submitted</p>
              <p class="thin">
                Supplying
                <span class="thick">{{ _fromAmount }} {{ _fromToken }}</span>
                and
                <span class="thick">{{ _toAmount }} {{ _toToken }}</span>
              </p>
              <br />
              <p class="sub">
                <a class="anchor" href="#"
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
import Loader from "@/components/shared/Loader.vue";
import SifButton from "@/components/shared/SifButton.vue";
export default defineComponent({
  components: { Loader, SifButton },
  props: {
    confirmed: Boolean,
    state: String,
    fromAmount: String,
    fromToken: String,
    toAmount: String,
    toToken: String,
  },
  setup(props) {
    // Need to cache amounts and disconnect reactivity
    return {
      _fromAmount: props.fromAmount,
      _fromToken: props.fromToken,
      _toAmount: props.toAmount,
      _toToken: props.toToken,
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
.thin {
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