<script lang="ts">
import { defineComponent, watch } from "vue";
import { computed, ref, ComputedRef } from "@vue/reactivity";
import Layout from "@/components/layout/Layout.vue";
import SifButton from "@/components/shared/SifButton.vue";
import {
  getAssetLabel,
  getBlockExplorerUrl,
  getRewardEarningsUrl,
  useAssetItem,
} from "@/components/shared/utils";
import { useCore } from "@/hooks/useCore";
import { useRoute } from "vue-router";
import { format } from "ui-core/src/utils/format";
import { Amount } from "ui-core";
import Tooltip from "@/components/shared/Tooltip.vue";
import Icon from "@/components/shared/Icon.vue";

const DECIMALS = 5;

async function getEarnedRewards(address: ComputedRef<any>, symbol: string) {
  const { config } = useCore();
  const earnedRewardsUrl = getRewardEarningsUrl(config.sifChainId);
  if (!address.value) return;
  const res = await fetch(
    `${earnedRewardsUrl}?symbol=${symbol}&address=${address.value}`,
  );
  const data = await res.json();
  const parsedData = JSON.parse(data);
  return parsedData.netChangeUSDT;
}

export default defineComponent({
  components: { Layout, SifButton, Tooltip, Icon },
  setup(props) {
    const { config, store } = useCore();
    const route = useRoute().params.externalAsset;

    const address = computed(() => store.wallet.sif.address);
    let earnedRewards = ref<string>("0");

    const accountPool = computed(() => {
      if (
        !route ||
        !store.wallet.sif.address ||
        !store.accountpools ||
        !store.accountpools[store.wallet.sif.address] ||
        !store.accountpools[store.wallet.sif.address][`${route}_rowan`]
      ) {
        return null;
      }

      const poolTicker = `${route}_rowan`;
      const storeAccountPool =
        store.accountpools[store.wallet.sif.address][poolTicker];

      // enrich pool ticker with pool object
      return {
        ...storeAccountPool,
        pool: store.pools[poolTicker],
      };
    });

    const fromSymbol = computed(() =>
      accountPool?.value?.pool.amounts[1].asset
        ? getAssetLabel(accountPool?.value.pool.amounts[1].asset)
        : "",
    );
    const fromAsset = useAssetItem(fromSymbol);
    const fromToken = fromAsset.token;
    const fromBackgroundStyle = fromAsset.background;
    const fromTokenImage = computed(() => {
      if (!fromToken.value) return "";
      const t = fromToken.value;
      return t.imageUrl;
    });

    watch([address, fromSymbol], async () => {
      const realEarnedRewards = await getEarnedRewards(
        address,
        fromSymbol.value?.toLowerCase() || "0",
      );
      earnedRewards.value = format(Amount(realEarnedRewards.toString()), {
        mantissa: 2,
      });
    });

    const fromTotalValue = computed(() => {
      const aAmount = accountPool?.value?.pool.amounts[1];
      if (!aAmount) return "";
      return format(aAmount.amount, aAmount.asset, { mantissa: DECIMALS });
    });

    const toSymbol = computed(() =>
      accountPool?.value?.pool.amounts[0].asset
        ? getAssetLabel(accountPool.value.pool.amounts[0].asset)
        : "",
    );
    const toAsset = useAssetItem(toSymbol);
    const toToken = toAsset.token;
    const toBackgroundStyle = toAsset.background;
    const toTokenImage = computed(() => {
      if (!toToken.value) return "";
      const t = toToken.value;
      return t.imageUrl;
    });

    const toTotalValue = computed(() => {
      const aAmount = accountPool?.value?.pool.amounts[0];
      if (!aAmount) return "";
      return format(aAmount.amount, aAmount.asset, { mantissa: DECIMALS });
    });

    const poolUnitsAsFraction = computed(
      () => accountPool?.value?.lp.units || Amount("0"),
    );

    const myPoolShare = computed(() => {
      if (!accountPool?.value?.pool?.poolUnits) return null;

      const perc = format(
        poolUnitsAsFraction.value
          .divide(accountPool?.value?.pool?.poolUnits)
          .multiply("100"),
        { mantissa: 4 },
      );

      return `${perc} %`;
    });
    const myPoolUnits = computed(() => {
      return format(poolUnitsAsFraction.value, { mantissa: DECIMALS });
    });
    return {
      accountPool,
      fromToken,
      fromSymbol,
      fromBackgroundStyle,
      fromTokenImage,
      fromTotalValue,
      toSymbol,
      toBackgroundStyle,
      toTokenImage,
      toTotalValue,
      myPoolUnits,
      myPoolShare,
      chainId: config.sifChainId,
      getBlockExplorerUrl,
      earnedRewards,
    };
  },
});
</script>

<template>
  <Layout class="pool" backLink="/pool" title="Your Pair">
    <div class="sheet" :class="!accountPool ? 'disabled' : 'active'">
      <div class="section">
        <div class="header" @click="$emit('poolselected')">
          <div class="image">
            <img
              v-if="fromTokenImage"
              width="22"
              height="22"
              :src="fromTokenImage"
              class="info-img"
            />
            <div class="placeholder" :style="fromBackgroundStyle" v-else></div>
            <img
              v-if="toTokenImage"
              width="22"
              height="22"
              :src="toTokenImage"
              class="info-img"
            />
            <div class="placeholder" :style="toBackgroundStyle" v-else></div>
          </div>
          <div class="symbol">
            <span>{{ fromSymbol }}</span>
            /
            <span>{{ toSymbol }}</span>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="details">
          <div
            class="row"
            :data-handle="'total-pooled-' + fromSymbol.toLowerCase()"
          >
            <span>Total Pooled {{ fromSymbol }}:</span>
            <span class="value">
              <span>{{ fromTotalValue }}</span>
              <img
                v-if="fromTokenImage"
                width="22"
                height="22"
                :src="fromTokenImage"
                class="info-img"
              />
              <div
                class="placeholder"
                :style="fromBackgroundStyle"
                v-else
              ></div>
            </span>
          </div>
          <div
            class="row"
            :data-handle="'total-pooled-' + toSymbol.toLowerCase()"
          >
            <span>Total Pooled {{ toSymbol.toUpperCase() }}:</span>
            <span class="value">
              <span>{{ toTotalValue }}</span>
              <img
                v-if="toTokenImage"
                width="22"
                height="22"
                :src="toTokenImage"
                class="info-img"
              />
              <div class="placeholder" :style="toBackgroundStyle" v-else></div>
            </span>
          </div>
          <div class="row" data-handle="total-pool-share">
            <span>Your pool share:</span>
            <span class="value">{{ myPoolShare }}</span>
          </div>
          <div class="row" data-handle="total-pool-share">
            <span
              >Your Net Gain/Loss:
              <Tooltip
                message="This is your net gain/loss based on earnings from swap fees and any gains or losses associated with changes in the tokens' prices. This is in USDT"
              >
                <Icon icon="info-box-black" /> </Tooltip
            ></span>
            <span class="value">${{ earnedRewards }}</span>
          </div>
        </div>
      </div>
      <div class="section">
        <div class="info">
          <h3 class="mb-2">Liquidity provider rewards</h3>
          <p class="text--small mb-2">
            Liquidity providers earn a percentage fee on all trades proportional
            to their share of the pool. Fees are added to the pool, accrue in
            real time and can be claimed by withdrawing your liquidity. To learn
            more, refer to the documentation
            <a
              target="_blank"
              href="https://docs.sifchain.finance/core-concepts/liquidity-pool"
              >here</a
            >.
          </p>
        </div>
      </div>

      <div class="section footer">
        <div class="mr-1">
          <router-link
            :to="`/pool/remove-liquidity/${fromSymbol.toLowerCase()}`"
            ><SifButton primaryOutline nocase block
              >Remove Liquidity</SifButton
            ></router-link
          >
        </div>
        <div class="ml-1">
          <router-link :to="`/pool/add-liquidity/${fromSymbol.toLowerCase()}`"
            ><SifButton primary nocase block
              >Add Liquidity</SifButton
            ></router-link
          >
        </div>
      </div>
      <div class="blockexplorer-container">
        <div class="blockexplorer-label">Blockexplorer</div>
        <div class="blockexplorer-link">
          <a target="_blank" :href="getBlockExplorerUrl(chainId)">View</a>
        </div>
      </div>
    </div>
  </Layout>
</template>

<style lang="scss" scoped>
.sheet {
  background: $c_white;
  border-radius: $br_sm;
  border: $divider;
  &.disabled {
    opacity: 0.3;
  }
  .section {
    padding: 8px 12px;
  }

  .section:not(:last-of-type) {
    border-bottom: $divider;
  }

  .header {
    display: flex;
  }
  .symbol {
    font-size: $fs_md;
    color: $c_text;
  }

  .image {
    height: 22px;

    & > * {
      border-radius: 16px;

      &:nth-child(2) {
        position: relative;
        left: -6px;
      }
    }
  }

  .row {
    display: flex;
    justify-content: space-between;
    padding: 2px 0;
    color: $c_text;
    font-weight: 400;

    .value {
      display: flex;
      align-items: center;
      font-weight: 700;
      & > * {
        margin-right: 0.5rem;
      }

      & > *:last-child {
        margin-right: 0;
      }
    }

    .image,
    .placeholder {
      margin-left: 4px;
    }
  }

  .info {
    text-align: left;
    font-weight: 400;
  }

  .placeholder {
    display: inline-block;
    background: #aaa;
    box-sizing: border-box;
    border-radius: 16px;
    height: 22px;
    width: 22px;
    text-align: center;
  }

  .footer {
    display: flex;

    & > div {
      flex: 1;
    }
  }
}
.blockexplorer-container {
  // TODO - This should be somewhat like the <Panel> class
  margin-top: 15px;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #dedede;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  .blockexplorer-label {
    color: #333;
    font-style: italic;
    font-size: 16px;
  }
  .blockexplorer-link {
    a {
      color: #666;
    }
    font-style: italic;
  }
}
</style>
