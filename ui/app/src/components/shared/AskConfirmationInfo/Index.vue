<script lang="tsx">
import AssetItemLarge from "@/components/shared/AssetItemLarge/AssetItemLarge.vue";
import ArrowIconButton from "@/components/shared/ArrowIconButton.vue";
import { defineComponent, useCssModule } from "vue";

export default defineComponent({
  props: {
    tokenAAmount: { type: String, required: true },
    tokenAToken: { type: String, required: true },
    tokenBAmount: { type: String, required: true },
    tokenBToken: { type: String, required: true },
  },
  setup(props) {
    const styles = useCssModule();
    return () => (
      <div class={styles.info}>
        <AssetItemLarge
          amount={props.tokenAAmount}
          symbol={props.tokenAToken}
        />
        <div class={styles.arrowHolder}>
          <ArrowIconButton enabled={true} />
        </div>
        <AssetItemLarge
          amount={props.tokenBAmount}
          symbol={props.tokenBToken}
        />
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
