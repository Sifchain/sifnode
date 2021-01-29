<script lang="ts">
import { PropType, defineComponent, ref } from "vue";
import { computed } from "@vue/reactivity";
import ConfirmationModal, {
  ConfirmState,
} from "@/components/shared/ConfirmationModal.vue";
import DetailsPanelPool from "@/components/shared/DetailsPanelPool.vue";
import AssetItemLarge from "@/components/shared/AssetItemLarge.vue";
import { useAssetItem } from "@/components/shared/utils";

export default defineComponent({
  components: {
    ConfirmationModal,
    DetailsPanelPool,
    AssetItemLarge,
  },

  props: {
    isOpen: { type: Boolean, default: false }, 
    requestClose: Function,
    state: { type: String as PropType<ConfirmState>, default: "confirming" },
    transactionHash: String,
    transactionStateMsg: String,
    fromToken: String,
    fromAmount: String,
    poolUnits: String,
    toAmount: String,
    toToken: String,
    aPerB: String,
    bPerA: String,
    shareOfPool: String,
  },

  setup(props, { emit }) {
    const fromSymbol = computed(() => props.fromToken);
    const fromAsset = useAssetItem(fromSymbol);
    const fromLabel = fromAsset.label;
    const fromImage = computed(() => fromAsset.token.value ? fromAsset.token.value.imageUrl : '');
    const fromAmount = computed(() => props.fromAmount);

    const toSymbol = computed(() => props.toToken);
    const toAsset = useAssetItem(toSymbol);
    const toLabel = toAsset.label;
    const toImage = computed(() => toAsset.token.value ? toAsset.token.value.imageUrl : '');
    const toAmount = computed(() => props.toAmount);

    // on confirmed, disconnect reactivity
    const _fromAmount = ref<string | undefined>("");
    const _fromToken = ref<string | undefined>("");
    const _toAmount = ref<string | undefined>("");
    const _toToken = ref<string | undefined>("");

    const onConfirmed = () => {
      _fromAmount.value = fromAmount.value;
      _fromToken.value = fromLabel.value;
      _toAmount.value = toAmount.value;
      _toToken.value = toLabel.value;
      emit('confirmed');
    }

    return {
      fromLabel,
      fromImage,
      fromAmount,
      toLabel,
      toImage,
      toAmount,

      _fromAmount,
      _fromToken,
      _toAmount,
      _toToken,

      onConfirmed,
    }
  },
});
</script>

<template>
  <ConfirmationModal 
    @confirmed="onConfirmed"
    :requestClose="requestClose"
    :isOpen="isOpen"
    :state="state"
    :transactionHash="transactionHash"
    :transactionStateMsg="transactionStateMsg"
    confirmButtonText="Confirm Supply"
    title="You are depositing"
  >

    <!-- Template for  -->
    <template v-slot:askbody>
      <div>
        <div class="pool-token-value mb-10">
          {{ poolUnits }}
        </div>
        <div class="pool-token">
          <div class="pool-token-image">
            <img
              v-if="fromImage"
              width="24"
              :src="fromImage"
              class="info-img"
            />
            <img
              v-if="toImage"
              width="24"
              :src="toImage"
              class="info-img"
            />
          </div>
          <div class="pool-token-label">
            {{ fromLabel }}/{{ toLabel }} Pool Tokens<br />
          </div>
        </div>
        <DetailsPanelPool
          class="details"
          :fromTokenLabel="fromLabel"
          :fromAmount="fromAmount"
          :fromTokenImage="fromImage"
          :toTokenLabel="toLabel"
          :toAmount="toAmount"
          :toTokenImage="toImage"
          :aPerB="aPerB"
          :bPerA="bPerA"
          :shareOfPool="shareOfPool"
        />
      </div>
    </template>

    <template v-slot:signing>
      <div>
        <p class="text--normal">
          Supplying
          <span class="text--bold">{{ _fromAmount }} {{ _fromToken }}</span>
          and
          <span class="text--bold">{{ _toAmount }} {{ _toToken }}</span>
        </p>
      </div>
    </template>

    <template v-slot:confirmed>
      <div>
        <p class="text--normal">
          Supplied
          <span class="text--bold">{{ _fromAmount }} {{ _fromToken }}</span>
          and
          <span class="text--bold">{{ _toAmount }} {{ _toToken }}</span>
        </p>
      </div>
    </template>

    <template v-slot:rejected>
      <div>
        <p class="text--normal">
          Supplying
          <span class="text--bold">{{ _fromAmount }} {{ _fromToken }}</span>
          and
          <span class="text--bold">{{ _toAmount }} {{ _toToken }}</span>
        </p>
      </div>
    </template>

    <template v-slot:failed>
      <div>
        <p class="text--normal">
          Supplying
          <span class="text--bold">{{ _fromAmount }} {{ _fromToken }}</span>
          and
          <span class="text--bold">{{ _toAmount }} {{ _toToken }}</span>
        </p>
      </div>
    </template>

  </ConfirmationModal>
</template>


<style lang="scss" scoped>
.details {
  margin-bottom: 20px;
  margin-top: 40px;
}

.pool-token {
  display: flex;
  margin-bottom: 8px;

  &-value {
    font-size: 30px;
    font-style: normal !important;
    text-align: left;
  }
  &-image {
    height: 26px;

    & > * {
      border-radius: 16px;

      &:nth-child(2) {
        position: relative;
        left: -8px;
      }
    }
  }
  &-label {
    text-align: left;
    font-weight: 400;
  }
  .placeholder {
    display: inline-block;
    background: #aaa;
    box-sizing: border-box;
    border-radius: 16px;
    height: 24px;
    width: 24px;
    text-align: center;
  }
}
</style>