<script lang="tsx">
import { defineComponent, PropType } from "vue";
import AskConfirmation from "./AskConfirmation.vue";
import AnimatedConfirmation from "./AnimatedConfirmation.vue";
import { IAssetAmount, TransactionStatus } from "ui-core";
import { effect } from "@vue/reactivity";

export default defineComponent({
  components: { AskConfirmation, AnimatedConfirmation },
  inheritAttrs: false,
  props: {
    state: {
      type: String as PropType<"confirm" | "submit" | "fail" | "success">,
      default: "confirm",
    },
    txStatus: { type: Object as PropType<TransactionStatus>, default: null },
    requestClose: Function as PropType<() => void>,
    priceMessage: { type: String, default: "" },
    fromAmount: { type: Object as PropType<IAssetAmount>, required: true },
    toAmount: { type: Object as PropType<IAssetAmount>, required: true },
    leastAmount: String,
    swapRate: String,
    minimumReceived: String,
    providerFee: String,
    priceImpact: String,
  },
  emits: ["confirmswap"],
  setup(props, ctx) {
    return () => (
      <>
        {props.state === "confirm" && (
          <AskConfirmation
            fromAmount={props.fromAmount}
            fromToken={props.fromAmount.label}
            toAmount={props.toAmount}
            toToken={props.toAmount.label}
            leastAmount={props.leastAmount}
            swapRate={props.swapRate}
            minimumReceived={props.minimumReceived}
            providerFee={props.providerFee}
            priceImpact={props.priceImpact}
            priceMessage={props.priceMessage}
            onConfirmswap={() => ctx.emit("confirmswap")}
          />
        )}
        {props.state === "submit" ||
          props.state === "fail" ||
          (props.state === "success" && (
            <AnimatedConfirmation
              state={props.state}
              txStatus={props.txStatus}
              fromAmount={props.fromAmount}
              toAmount={props.toAmount}
              onCloserequested={props.requestClose}
            />
          ))}
      </>
    );
  },
});
</script>
