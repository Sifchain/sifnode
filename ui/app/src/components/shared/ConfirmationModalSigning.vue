<template>
  <div>
    <div class="confirmation">
      <div class="message">
        <Loader black :success="confirmed" /><br />
        <div class="text-wrapper">
          <transition name="swipe">
            <div class="text" v-if="state === 'signing'">
              <slot name="signing"></slot>
            </div>
          </transition>
          <transition name="swipe">
            <div class="text" v-if="confirmed">
              <slot name="success"></slot>
            </div>
          </transition>
          <transition name="swipe">
            <div class="text" v-if="failed">
              <slot name="error"></slot>
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
import Loader from "@/components/shared/Loader.vue";
import SifButton from "@/components/shared/SifButton.vue";

export default {
  components: { Loader, SifButton },
  props: {
    confirmed: Boolean,
    failed: Boolean,
    state: String,
  },
};
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