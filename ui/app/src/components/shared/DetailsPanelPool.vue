<style lang="scss" module>
.details {
  border: 1px solid $c_gray_200;
  border-radius: $br_sm;
  background: $c_white;
}
.detailsHeader {
  padding: 10px 15px;
  border-bottom: 1px solid $c_gray_200;
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
import { format, IAssetAmount } from "ui-core";

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
    const { aPerB, bPerA } = props;
    const realAPerB = computed(() => {
      return aPerB === "N/A" ? "0" : aPerB;
    });
    const realBPerA = computed(() => {
      return bPerA === "N/A" ? "0" : bPerA;
    });

    return () => (
      <div class={styles.details}>
        <div class={styles.detailsHeader}>
          <div
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
          </div>
        </div>
        <div class={styles.detailsBody}>
          {realBPerA && (
            <div class={styles.detailsRow} data-handle="real-b-per-a-row">
              <span>Rates</span>
              <span>
                1
                {props.tokenAAmount.symbol.toLowerCase().includes("rowan")
                  ? props.tokenAAmount.symbol.toUpperCase()
                  : "c" + props.tokenAAmount.symbol.slice(1).toUpperCase()}
                = {realBPerA}
                {props.tokenBAmount.symbol.toLowerCase().includes("rowan")
                  ? props.tokenBAmount.symbol.toUpperCase()
                  : "c" + props.tokenBAmount.symbol.slice(1).toUpperCase()}
              </span>
            </div>
          )}
          {realAPerB && (
            <div class={styles.detailsRow} data-handle="real-a-per-b-row">
              <span>&nbsp;</span>
              <span>
                1
                {props.tokenBAmount.symbol.toLowerCase().includes("rowan")
                  ? props.tokenBAmount.symbol.toUpperCase()
                  : "c" + props.tokenBAmount.symbol.slice(1).toUpperCase()}
                = {realAPerB}
                {props.tokenAAmount.symbol.toLowerCase().includes("rowan")
                  ? props.tokenAAmount.symbol.toUpperCase()
                  : "c" + props.tokenAAmount.symbol.slice(1).toUpperCase()}
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
