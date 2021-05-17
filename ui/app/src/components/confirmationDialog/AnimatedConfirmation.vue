<script lang="tsx">
import { defineComponent, PropType, useCssModule } from "vue";
import Loader from "@/components/shared/Loader.vue";
import SifButton from "@/components/shared/SifButton.vue";
import { useCore } from "@/hooks/useCore";
import { getBlockExplorerUrl } from "../shared/utils";
import { ErrorCode, TransactionStatus } from "ui-core";
import SwipeTransition from "./SwipeTransition.vue";

export default defineComponent({
  components: { Loader, SifButton },
  props: {
    txStatus: { type: Object as PropType<TransactionStatus>, default: null },
    confirmed: Boolean,
    failed: Boolean,
    state: { type: String as PropType<"submit" | "fail" | "success"> },
    fromAmount: String,
    fromToken: String,
    toAmount: String,
    toToken: String,
  },
  setup(props, context) {
    const { config } = useCore();
    const styles = useCssModule();

    // Create a template for our confirmation screens
    const ConfirmTemplate = (p: {
      header: JSX.Element;
      pre: JSX.Element;
      fromAmount?: string;
      toAmount?: string;
      fromToken?: string;
      toToken?: string;
      post: JSX.Element;
    }) => (
      <div class={styles.text}>
        <p>{p.header}</p>
        <p class={styles.thin} data-handle="swap-message">
          {p.pre + " "}
          <span class={styles.thick}>
            {p.fromAmount} {p.fromToken}
          </span>{" "}
          for{" "}
          <span class={styles.thick}>
            {p.toAmount} {p.toToken}
          </span>
        </p>
        <br />
        <p class={styles.sub}>{p.post}</p>
      </div>
    );

    // Need to cache amounts and disconnect reactivity
    const amounts = {
      fromAmount: props.fromAmount,
      fromToken: props.fromToken,
      toAmount: props.toAmount,
      toToken: props.toToken,
    };

    const handleCloseRequested = () => context.emit("closerequested");

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
                {props.state === "fail" &&
                  props.txStatus.code === ErrorCode.TX_FAILED_OUT_OF_GAS && (
                    <ConfirmTemplate
                      header="Transaction Failed"
                      pre="Failed to swap"
                      post="Please try to increase the gas limit."
                      {...amounts}
                    />
                  )}
              </SwipeTransition>
              <SwipeTransition>
                {props.state === "fail" &&
                  props.txStatus.code === ErrorCode.USER_REJECTED && (
                    <ConfirmTemplate
                      header="Transaction Rejected"
                      pre="Failed to swap"
                      post="Please confirm the transaction in your wallet."
                      {...amounts}
                    />
                  )}
              </SwipeTransition>
              <SwipeTransition>
                {props.state === "fail" &&
                  props.txStatus.code === ErrorCode.TX_FAILED_SLIPPAGE && (
                    <ConfirmTemplate
                      header="Transaction Failed"
                      pre="Failed to swap"
                      post="Please try to increase slippage tolerance."
                      {...amounts}
                    />
                  )}
              </SwipeTransition>
              <SwipeTransition>
                {props.state === "fail" && (
                  <ConfirmTemplate
                    header="Transaction Failed"
                    pre="Failed to swap"
                    post=""
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
