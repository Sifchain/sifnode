<script lang="tsx">
import { computed, defineComponent, PropType, useCssModule } from "vue";
import Icon from "@/components/shared/Icon.vue";
import Tooltip from "@/components/shared/Tooltip.vue";
import { getMantissaLength, format, IAssetAmount } from "ui-core";
import { useAssetItem } from "../utils";
import { effect } from "@vue/reactivity";

export default defineComponent({
  props: {
    amount: { type: Object as PropType<IAssetAmount>, required: true },
  },
  setup(props) {
    effect(() => {
      console.log({ "props.amount": props.amount });
    });
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
        <div class={styles.infoAmount} data-handle="info-amount">
          {format(props.amount.amount, props.amount.asset, { mantissa: 6 })}
        </div>
        {props.amount && getMantissaLength(props.amount.toDerived()) > 6 ? (
          <Tooltip message={props.amount.toString()} fit>
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
.infoToken {
  font-size: 18px;
  width: 60px;
}
.infoImg {
  width: 20px;
  margin-right: 10px;
}
.infoEye {
  width: 30px;
}
</style>
