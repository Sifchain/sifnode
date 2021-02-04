<template>
  <div class="positioner">
    <LoaderCircle :black="black" :success="success || failed" />
    <div class="tick-holder">
      <transition name="reveal">
        <LoaderTick v-if="success" class="tick" />
      </transition>
      <transition name="reveal">
        <LoaderFailed v-if="failed" class="tick" />
      </transition>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import LoaderTick from "@/components/shared/LoaderTick.vue";
import LoaderCircle from "@/components/shared/LoaderCircle.vue";
import LoaderFailed from "@/components/shared/LoaderFailed.vue";

export default defineComponent({
  components: { LoaderCircle, LoaderTick, LoaderFailed },
  props: {
    black: { type: Boolean, default: false },
    success: { type: Boolean, default: false },
    failed: { type: Boolean, default: false },
  },
});
</script>

<style lang="scss" scoped>
.positioner {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
}

.tick-holder {
  position: absolute;
  width: 66px;
  height: 66px;
}

.tick {
  position: absolute;
  overflow: hidden;
  width: 66px;
  height: 66px;
  stroke: $c_gray_700;
  max-width: 66px;
}
.reveal-leave-active {
  transition: max-width 0.01s ease-out;
}
.reveal-enter-active {
  transition: max-width 1s ease-out;
}

.reveal-enter-from {
  max-width: 0;
}
</style>
