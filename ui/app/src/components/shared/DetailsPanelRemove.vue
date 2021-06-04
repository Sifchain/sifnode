<template>
  <div>
    <Copy class="mb-8">
      Are you sure you want to withdraw your liquidity? This withdrawal will not
      automatically submit a claim for the associated rewards. While it will
      reset your reward tickets, you will still need to go into the DEX and
      manually submit a claim transaction to fully claim those rewards.
    </Copy>
    <Copy class="mb-8">
      Unless you have reached full maturity on your rewards, we recommend that
      you do not withdraw your liquidity so you can realize your full rewards.
    </Copy>
    <Copy class="mb-8">
      Find out
      <a
        href="https://docs.sifchain.finance/resources/rewards-programs"
        target="_blank"
        >more information here</a
      >.
    </Copy>
    <PairTable class="mb-8" :items="computedPairPanel" />
    <div class="details mb-8">
      <div class="details-header">
        <h4 class="text--left">You Should Receive:</h4>

        <div class="details-row">
          <AssetItem :symbol="nativeAssetSymbol" /> {{ nativeAssetAmount }}
        </div>
        <div class="details-row">
          <AssetItem :symbol="externalAssetSymbol" /> {{ externalAssetAmount }}
        </div>
      </div>
    </div>
    <br /><br />
  </div>
</template>

<style lang="scss" scoped>
.details {
  border: 1px solid $c_gray_200;
  border-radius: $br_sm;
  background: $c_white;

  &-header {
    padding: 10px 15px;
    border-bottom: 1px solid $c_gray_200;
  }
  &-body {
    padding: 10px 15px;
  }

  &-row {
    display: flex;
    flex-direction: row;
    align-items: center;

    span:last-child {
      text-align: right;
      color: $c_gray_900;
      margin-left: auto;
    }

    span:first-child {
      color: $c_gray_700;
      font-weight: 400;
      text-align: left;
    }

    img {
      margin-right: 8px;
    }
  }
}
</style>
<script lang="ts">
import { defineComponent } from "vue";
import AssetItem from "@/components/shared/AssetItem.vue";
import { Copy } from "@/components/shared/Text";
import PairTable from "@/components/shared/PairTable.vue";

export default defineComponent({
  components: {
    AssetItem,
    Copy,
    PairTable,
  },
  props: {
    rewardsData: { type: Object, default: {} },
    nativeAssetSymbol: { type: String, default: "" },
    nativeAssetAmount: { type: String, default: "" },
    externalAssetAmount: { type: String, default: "" },
    externalAssetSymbol: { type: String, default: "" },
  },
  computed: {
    computedPairPanel(): Array<Object> {
      return [
        {
          key: "Claimable  Rewards",
          value: this.$props.rewardsData
            .totalClaimableCommissionsAndClaimableRewards,
        },
        {
          key: "Projected Full Amount",
          value: this.$props.rewardsData.totalCommissionsAndRewardsAtMaturity,
        },
        {
          key: "Maturity Date",
          value: this.$props.rewardsData.maturityDateISO,
          type: "date",
        },
      ];
    },
  },
});
</script>
