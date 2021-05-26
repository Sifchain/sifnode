<script lang="tsx">
import AssetItemLarge from "@/components/shared/AssetItemLarge/AssetItemLarge.vue";
import ArrowIconButton from "@/components/shared/ArrowIconButton.vue";
import { defineComponent, PropType, useCssModule } from "vue";
import { IAssetAmount } from "ui-core";
import { effect } from "@vue/reactivity";

export default defineComponent({
  props: {
    tokenAAmount: { type: Object as PropType<IAssetAmount>, required: true },
    tokenBAmount: { type: Object as PropType<IAssetAmount>, required: true },
  },
  setup(props) {
    const styles = useCssModule();

    return () => (
      <div class={styles.info}>
        <AssetItemLarge amount={props.tokenAAmount} />
        <div class={styles.arrowHolder}>
          <ArrowIconButton enabled={true} />
        </div>
        <AssetItemLarge amount={props.tokenBAmount} />
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
</style>
