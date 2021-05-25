<script lang="tsx">
import { defineComponent, PropType } from "vue";
import AskConfirmation from "./AskConfirmation.vue";
import AnimatedConfirmation from "./AnimatedConfirmation.vue";
import { IAssetAmount, TransactionStatus } from "ui-core";

export default defineComponent({
  components: { AskConfirmation, AnimatedConfirmation },
  inheritAttrs: false,
  props: {
    state: {
      type: String as PropType<"confirm" | "submit" | "fail" | "success">,
      default: "confirm",
    },
    txStatus: { type: Object as PropType<TransactionStatus>, default: null },
    requestClose: Function,
    priceMessage: { type: String, default: "" },
    fromAmount: { type: Object as PropType<IAssetAmount>, required: true },
    toAmount: { type: Object as PropType<IAssetAmount> },
    toToken: String,
    leastAmount: String,
    swapRate: String,
    minimumReceived: String,
    providerFee: String,
    priceImpact: String,
  },
  emits: ["confirmswap"],
  setup(props, ctx) {
    return (
      <>
        {props.state === "confirm" && (
          <AskConfirmation
            fromAmount={props.fromAmount}
            fromToken={props.fromAmount.label}
            toAmount="toAmount"
            toToken="toToken"
            leastAmount="leastAmount"
            swapRate="swapRate"
            minimumReceived="minimumReceived"
            providerFee="providerFee"
            priceImpact="priceImpact"
            priceMessage="priceMessage"
            onConfirmswap={ctx.emit("confirmswap")}
          />
        )}
        {props.state === "submit" ||
          props.state === "fail" ||
          (props.state === "success" && (
            <AnimatedConfirmation
              state={props.state}
              txStatus={props.txStatus}
              fromAmount={props.fromAmount.toString()}
              fromToken={props.fromAmount.label}
              toAmount="toAmount"
              toToken="toToken"
              on-closerequested={props.requestClose}
            />
          ))}
      </>
    );
  },
});
</script>
