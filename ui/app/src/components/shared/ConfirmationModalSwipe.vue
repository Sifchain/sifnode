<script lang="tsx">
import { computed, defineComponent, PropType, Transition } from "vue";
import { useCore } from "@/hooks/useCore";
// import Loader from "@/components/shared/Loader.vue";
// import SifButton from "@/components/shared/SifButton.vue";
import ConfirmationModalSwipeOuter from "./ConfirmationModalSwipeOuter.vue";
import ConfirmationModalSwipePanel from "./ConfirmationModalSwipePanel.vue";

export default defineComponent({
  props: {
    state: String,
    loaderState: Object as PropType<{
      [s: string]: { success: boolean; failed: boolean };
    }>,
  },

  setup(props, context) {
    // const { config } = useCore();
    const stateMap = computed(() => {
      if (!props.loaderState || !props.state) return null;

      return props.loaderState[props.state];
    });
    const success = computed(() => {
      if (!stateMap.value) return false;

      return stateMap.value.success || false;
    });

    const failed = computed(() => {
      if (!stateMap.value) return false;

      return stateMap.value.failed || false;
    });

    return () => {
      return (
        <ConfirmationModalSwipeOuter
          success={success.value}
          failed={failed.value}
        >
          {Object.entries(context.slots).map(([name, slotFn]) => {
            if (!slotFn) return null;
            return (
              <Transition name="swipe">
                <ConfirmationModalSwipePanel v-show={props.state === name}>
                  {slotFn()}
                </ConfirmationModalSwipePanel>
              </Transition>
            );
          })}
        </ConfirmationModalSwipeOuter>
      );
    };
  },
});
</script>

<style lang="scss">
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