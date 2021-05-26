<style lang="scss" module>
.detailsHeader {
  padding: 10px 15px;
  // border-bottom: 1px solid $c_gray_200;
}
.detailsBody {
  padding: 10px 15px;
}

.detailsRow {
  display: flex;
  justify-content: space-between;

  span:last-child {
    text-align: right;
    color: $c_gray_900;
  }

  > span:first-child {
    color: $c_gray_700;
    font-weight: 400;
    text-align: left;
  }
}
.detailsRowAsset {
  display: flex;
  align-items: center;
}

.detailsRowValue {
  display: flex;
  color: $c_black;
  img {
    margin-left: 5px;
  }
}
</style>
<script lang="tsx">
import { defineComponent, PropType, useCssModule } from "vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import { computed } from "@vue/reactivity";
import { IAssetAmount } from "ui-core";
import AskConfirmationInfo from "@/components/shared/AskConfirmationInfo/Index.vue";

export default defineComponent({
  components: {
    AssetItem,
  },
  props: {
    tokenAAmount: { type: Object as PropType<IAssetAmount>, default: null },
    tokenBAmount: { type: Object as PropType<IAssetAmount>, default: null },
    aPerB: { type: String, default: "" },
    bPerA: { type: String, default: "" },
    shareOfPool: String,
  },
  setup(props) {
    const styles = useCssModule();

    const realAPerB = computed(() => {
      return props.aPerB === "N/A" ? "0" : props.aPerB;
    });
    const realBPerA = computed(() => {
      return props.bPerA === "N/A" ? "0" : props.bPerA;
    });

    return () => (
      <div>
        <AskConfirmationInfo
          amountDescriptions="Deposited"
          tokenAAmount={props.tokenAAmount}
          tokenBAmount={props.tokenBAmount}
        />
        <div class={styles.detailsHeader}>
          {/* <div
            class={styles.detailsRow}
            data-handle="token-a-details-panel-pool-row"
          >
            <span class={styles.detailsRowAsset}>
              {props.tokenAAmount && (
                <AssetItem symbol={props.tokenAAmount.symbol} inline />
              )}
              &nbsp;Deposited
            </span>
            <div class={styles.detailsRowValue} data-handle="details-row-value">
              <span>
                {format(props.tokenAAmount.amount, props.tokenAAmount.asset, {
                  mantissa: 18,
                }) || "0"}
              </span>
            </div>
          </div>
          <div
            class={styles.detailsRow}
            data-handle="token-b-details-panel-pool-row"
          >
            <span class={styles.detailsRowAsset}>
              {props.tokenBAmount && (
                <AssetItem symbol={props.tokenBAmount.symbol} inline />
              )}
              &nbsp;Deposited
            </span>
            <div class={styles.detailsRowValue} data-handle="details-row-value">
              <span>
                {format(props.tokenBAmount.amount, props.tokenBAmount.asset, {
                  mantissa: 18,
                }) || 0}
              </span>
            </div>
          </div>*/}
        </div>
        <div class={styles.detailsBody}>
          {realBPerA.value && (
            <div class={styles.detailsRow} data-handle="real-b-per-a-row">
              <span>Rates</span>
              <span>
                {`1 ${props.tokenAAmount.label} = ${realBPerA.value} ${props.tokenBAmount.label}`}
              </span>
            </div>
          )}
          {realAPerB.value && (
            <div class={styles.detailsRow} data-handle="real-a-per-b-row">
              <span>&nbsp;</span>
              <span>
                {`1 ${props.tokenBAmount.label} = ${realAPerB.value} ${props.tokenAAmount.label}`}
              </span>
            </div>
          )}
          <div class={styles.detailsRow} data-handle="real-share-of-pool">
            <span>Share of Pool:</span>
            <span>{props.shareOfPool}</span>
          </div>
        </div>
      </div>
    );
  },
});
</script>
