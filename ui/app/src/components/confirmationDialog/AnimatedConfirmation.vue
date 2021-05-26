<script lang="tsx">
import { defineComponent, PropType, useCssModule } from "vue";
import Loader from "@/components/shared/Loader.vue";
import SifButton from "@/components/shared/SifButton.vue";
import { useCore } from "@/hooks/useCore";
import { getBlockExplorerUrl } from "../shared/utils";
import { ErrorCode, format, IAssetAmount, TransactionStatus } from "ui-core";
import SwipeTransition from "./SwipeTransition.vue";

export default defineComponent({
  components: { Loader, SifButton },
  emits: ["closerequested"],
  props: {
    txStatus: { type: Object as PropType<TransactionStatus>, default: null },
    confirmed: Boolean,
    failed: Boolean,
    state: { type: String as PropType<"submit" | "fail" | "success"> },
    fromAmount: { type: Object as PropType<IAssetAmount>, required: true },
    toAmount: { type: Object as PropType<IAssetAmount>, required: true },
    onCloserequested: Function as PropType<() => void>,
  },
  setup(props, context) {
    const { config } = useCore();
    const styles = useCssModule();

    // Create a template for our confirmation screens
    const ConfirmTemplate = (p: {
      header: JSX.Element;
      pre: JSX.Element;
      fromAmount: IAssetAmount;
      toAmount: IAssetAmount;
      post: JSX.Element;
    }) => (
      <div class={styles.text}>
        <p>{p.header}</p>
        <p class={styles.thin} data-handle="swap-message">
          {p.pre + " "}
          <span class={styles.thick}>
            {format(p.fromAmount.amount, p.fromAmount.asset, { mantissa: 6 })}{" "}
            {p.fromAmount.label}
          </span>{" "}
          for{" "}
          <span class={styles.thick}>
            {format(p.toAmount.amount, p.toAmount.asset, { mantissa: 6 })}{" "}
            {p.toAmount.label}
          </span>
        </p>
        <br />
        <p class={styles.sub}>{p.post}</p>
      </div>
    );

    // Need to cache amounts and disconnect reactivity
    const amounts = {
      fromAmount: props.fromAmount,
      toAmount: props.toAmount,
    };

    const handleCloseRequested = () => context.emit("closerequested");

    const parseErrorObject = (
      txStatus: TransactionStatus,
    ): { header: JSX.Element; pre: JSX.Element; post: JSX.Element } => {
      const messageLookup = {
        [ErrorCode.TX_FAILED_SLIPPAGE]:
          "Please try to increase slippage tolerance.",
        [ErrorCode.USER_REJECTED]: "You have rejected the transaction.",
        [ErrorCode.TX_FAILED_OUT_OF_GAS]:
          "Please try to increase the gas limit.",
        [ErrorCode.TX_FAILED_NOT_ENOUGH_ROWAN_TO_COVER_GAS]:
          "Not enough ROWAN to cover the gas fees. Please try again and ensure you have enough ROWAN to cover the selected gas fees.",
      };

      const post =
        typeof txStatus.code !== "undefined"
          ? messageLookup[txStatus.code as keyof typeof messageLookup] ||
            txStatus.memo ||
            ""
          : "";

      return {
        header: "Transaction Failed",
        pre: "Failed to swap",
        post,
      };
    };
    return () => (
      <div>
        <div class={styles.confirmation}>
          <div class={styles.message}>
            <Loader
              black
              success={props.state === "success"}
              failed={props.state === "fail"}
            />
            <br />
            <div class={styles.textWrapper}>
              <SwipeTransition>
                {props.state === "submit" && (
                  <ConfirmTemplate
                    header="Waiting for confirmation"
                    pre="Swapping"
                    post="Confirm this transaction in your wallet"
                    {...amounts}
                  />
                )}
              </SwipeTransition>
              <SwipeTransition>
                {props.state === "fail" && (
                  <ConfirmTemplate
                    {...parseErrorObject(props.txStatus)}
                    {...amounts}
                  />
                )}
              </SwipeTransition>
              <SwipeTransition>
                {props.state === "success" && (
                  <ConfirmTemplate
                    header="Transation Submitted"
                    pre="Swapped"
                    post={
                      <a
                        class={styles.anchor}
                        target="_blank"
                        href={getBlockExplorerUrl(
                          config.sifChainId,
                          props.txStatus.hash,
                        )}
                      >
                        View transaction on Block Explorer
                      </a>
                    }
                    {...amounts}
                  />
                )}
              </SwipeTransition>
            </div>
          </div>
        </div>
        <div
          class={
            props.state === "success" ? styles.footerConfirmed : styles.footer
          }
        >
          <SifButton block {...{ onClick: handleCloseRequested }} primary>
            Close
          </SifButton>
        </div>
      </div>
    );
  },
});
</script>
<style lang="scss" module>
.confirmation {
  display: flex;
  justify-content: start;
  align-items: start;
  min-height: 40vh;
  padding: 15px 20px;
}
.message {
  width: 100%;
  font-size: 18px;
  margin-top: 3em;
}
.textWrapper {
  margin-top: 0.5em;
  position: relative;
  display: flex;
  width: 100%;
}
.text {
  position: absolute;
  width: 100%;
}
.anchor {
  color: $c_black;
}
.thin {
  margin-top: 1em;
  margin-bottom: 1em;
  font-size: 15px;
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
}
.footerConfirmed {
  padding: 16px;
  transition: opacity 0.5s ease-out;
  opacity: 1;
  visibility: inherit;
}
</style>
