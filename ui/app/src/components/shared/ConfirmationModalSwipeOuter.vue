<template>
  <div class="modal-inner">
    <div>
      <div class="confirmation">
        <div class="message">
          <Loader black :success="success" :failed="failed" />
          <br />
          <div class="text-wrapper">
            <slot name="default"></slot>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent } from "vue";
import Loader from "./Loader.vue";
export default defineComponent({
  components: {
    Loader,
  },
  props: {
    success: Boolean,
    failed: Boolean,
  },
});
</script>
<style lang="scss" scoped>
.modal-inner {
  display: flex;
  flex-direction: column;
  padding: 30px 20px 20px 20px;
  min-height: 50vh;
}
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