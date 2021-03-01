<script lang="ts">
import { computed, ref } from "@vue/reactivity";
import { defineComponent, PropType } from "vue";
import { getAssetLabel, useAssetItem } from "@/components/shared/utils";
import { LiquidityProvider, Pool } from "ui-core";
import { useRouter } from "vue-router";
import AssetItemPool from "@/components/shared/AssetItemPool.vue";

export default defineComponent({
  props: {
    accountPool: {
      type: Object as PropType<{ lp: LiquidityProvider; pool: Pool }>,
    },
  },

  components: {
    AssetItemPool,
  },

  setup(props) {
    const router = useRouter();

    const fromSymbol = computed(() =>
      props.accountPool?.pool.amounts[1].asset
        ? getAssetLabel(props.accountPool?.pool.amounts[1].asset)
        : ""
    );

    const handleClick = () => {
      router.push(`/pool/${fromSymbol.value.toLowerCase()}`);
    };

    return {
      handleClick,
    };
  },
});
</script>

<template>
  <div class="pool-list-item" @click="handleClick">
    <AssetItemPool :pool="accountPool" />
    <div class="button">Manage</div>
  </div>
</template>

<style scoped lang="scss">
.pool-list-item {
  padding: 14px 16px;
  display: flex;

  &:not(:last-of-type) {
    border-bottom: $divider;
  }

  &:hover {
    cursor: pointer;
    background: $c_gray_50;
  }

  .button {
    font-size: $fs_sm;
    font-weight: normal;
    flex-grow: 1;
    text-align: right;
  }
}
</style>
