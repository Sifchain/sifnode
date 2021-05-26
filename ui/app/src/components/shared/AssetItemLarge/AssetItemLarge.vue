<script lang="tsx">
import { computed, defineComponent, PropType, useCssModule } from "vue";
import Icon from "@/components/shared/Icon.vue";
import Tooltip from "@/components/shared/Tooltip.vue";
import { getMantissaLength, format, IAssetAmount } from "ui-core";
import { useAssetItem } from "../utils";

export default defineComponent({
  props: {
    amount: { type: Object as PropType<IAssetAmount>, required: true },
    description: { type: String, required: false },
  },
  setup(props) {
    const styles = useCssModule();
    const symbol = computed(() => props.amount?.symbol);
    const { token, label: tokenLabel } = useAssetItem(symbol);

    const tokenImage = computed(() => {
      if (!token.value) return "";
      const t = token.value;
      return t.imageUrl;
    });

    return () => (
      <div
        class={styles.infoRow}
        data-handle={"info-row-" + tokenLabel.value.toLowerCase()}
      >
        {tokenImage.value && (
          <img width="24" src={tokenImage.value} class={styles.infoImg} />
        )}
        <div class={styles.infoToken}>{tokenLabel.value}</div>
        {props.description && (
          <div class={styles.infoDescription}>{props.description}</div>
        )}
        <div class={styles.infoAmount} data-handle="info-amount">
          {format(props.amount.amount, props.amount.asset, { mantissa: 6 })}
        </div>
        {props.amount && getMantissaLength(props.amount.toDerived()) > 6 ? (
          <Tooltip
            message={
              format(props.amount.amount, props.amount.asset) +
              " " +
              props.amount.label
            }
            fit
          >
            <div class={styles.iconHolder}>
              <Icon icon="eye" class={styles.infoEye} />
            </div>
          </Tooltip>
        ) : (
          <div class={styles.eyePlaceholder} />
        )}
      </div>
    );
  },
});
</script>

<style lang="scss" module>
.eyePlaceholder {
  width: 21px;
}
.iconHolder {
  margin-left: 5px;
  margin-top: 7px;
  svg {
    fill: #c6c6c6 !important;
    width: 16px;
  }
  &:hover {
    svg {
      fill: #d4b553 !important;
    }
  }
}
.infoRow {
  display: flex;
  justify-content: start;
  align-items: center;
  font-weight: 400;
}
.infoAmount {
  font-size: 18px;
  flex: 1 1 auto;
  text-align: right;
}
.infoDescription {
  color: $c_gray_700;
  transform: translateY(2px);
  font-size: 14px;
  margin-left: 5px;
}
.infoToken {
  font-size: 18px;
  min-width: 68px;
  text-align: left;
}
.infoImg {
  width: 20px;
  margin-right: 10px;
}
.infoEye {
  width: 30px;
}
</style>
