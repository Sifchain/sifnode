<script lang="tsx">
import AssetItemLarge from "@/components/shared/AssetItemLarge/AssetItemLarge.vue";
import ArrowIconButton from "@/components/shared/ArrowIconButton.vue";
import { defineComponent, PropType, useCssModule } from "vue";
import { IAssetAmount } from "ui-core";

export default defineComponent({
  props: {
    amountDescriptions: { type: String },
    tokenAAmount: { type: Object as PropType<IAssetAmount>, required: true },
    tokenBAmount: { type: Object as PropType<IAssetAmount>, required: true },
    showArrow: { type: Boolean, default: true },
  },
  setup(props) {
    const styles = useCssModule();

    return () => (
      <div class={styles.info}>
        <div data-handle="token-a-row">
          <AssetItemLarge
            description={props.amountDescriptions}
            amount={props.tokenAAmount}
          />
        </div>
        <div
          class={props.showArrow ? styles.arrowHolder : styles.arrowHolderEmpty}
        >
          {props.showArrow && <ArrowIconButton enabled={true} />}
        </div>
        <div data-handle="token-b-row">
          <AssetItemLarge
            description={props.amountDescriptions}
            amount={props.tokenBAmount}
          />
        </div>
      </div>
    );
  },
});
</script>

<style lang="scss" module>
.info {
  background: $c_gray_100;
  padding: 20px 20px 20px 20px;
  border-radius: $br_sm;
}

.arrowHolder {
  padding: 5px 4px;
}

.arrowHolderEmpty {
  height: 17px;
}
</style>
