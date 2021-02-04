<script lang="tsx">
import { computed, defineComponent, PropType, Transition } from "vue";
import { useCore } from "@/hooks/useCore";
// import Loader from "@/components/shared/Loader.vue";
// import SifButton from "@/components/shared/SifButton.vue";
import ConfirmationModalSwipeOuter from "./ConfirmationModalSwipeOuter.vue";
import ConfirmationModalSwipePanel from "./ConfirmationModalSwipePanel.vue";
import SwipeTransition from "./SwipeTransition.vue";

type LoaderState = {
  [s: string]: { success: boolean; failed: boolean };
};

export default defineComponent({
  props: {
    state: { type: String, required: true },
    loaderState: {
      type: Object as PropType<LoaderState>,
      required: true,
    },
  },

  setup(props, context) {
    return () => {
      const stateMap = props.loaderState[props.state];

      return (
        <ConfirmationModalSwipeOuter
          success={!!stateMap?.success}
          failed={!!stateMap?.failed}
        >
          {Object.entries(context.slots).map(([name, slotFn]) => {
            if (!slotFn) return null;
            return (
              <SwipeTransition>
                <ConfirmationModalSwipePanel v-show={props.state === name}>
                  {slotFn()}
                </ConfirmationModalSwipePanel>
              </SwipeTransition>
            );
          })}
        </ConfirmationModalSwipeOuter>
      );
    };
  },
});
</script>
