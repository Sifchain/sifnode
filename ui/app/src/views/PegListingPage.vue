<script lang="ts">
import Tab from "@/components/shared/Tab.vue";
import Tabs from "@/components/shared/Tabs.vue";
import Layout from "@/components/layout/Layout.vue";
import AssetList from "@/components/shared/AssetList.vue";
import SifInput from "@/components/shared/SifInput.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import SifButton from "@/components/shared/SifButton.vue";
import Tooltip from "@/components/shared/Tooltip.vue";
import Icon from "@/components/shared/Icon.vue";

import { sortAssetAmount } from "./utils/sortAssetAmount";
import { useCore } from "@/hooks/useCore";
import { defineComponent, ref } from "vue";
import { computed } from "@vue/reactivity";
import { getUnpeggedSymbol } from "../components/shared/utils";
import { AssetAmount, IAsset, IAssetAmount, TransactionStatus } from "ui-core";
type TokenListItem = {
  amount: IAssetAmount;
  asset: IAsset;
  pegTxs: TransactionStatus[];
  supported: boolean;
};
export default defineComponent({
  components: {
    Tab,
    Tabs,
    AssetList,
    Layout,
    SifButton,
    SifInput,
    ActionsPanel,
    Tooltip,
    Icon,
  },
  setup(_, context) {
    const { store, actions } = useCore();
    function getIsSupportedNetwork(asset: IAsset): boolean {
      if (asset.network === "ethereum") {
        return actions.ethWallet.isSupportedNetwork();
      }

      if (asset.network === "sifchain") {
        return true; // TODO: Handle the case of whether the network is supported
      }
      return false;
    }
    const searchText = ref("");
    const selectedTab = ref("Sifchain Native");

    const allTokens = computed(() => {
      if (selectedTab.value === "External Tokens") {
        return actions.peg.getEthTokens();
      }

      if (selectedTab.value === "Sifchain Native") {
        return actions.peg.getSifTokens();
      }
      return [];
    });

    const pendingPegTxList = computed(() => {
      if (
        !store.wallet.eth.address ||
        !store.tx.eth ||
        !store.tx.eth[store.wallet.eth.address]
      )
        return null;

      const txs = store.tx.eth[store.wallet.eth.address];

      const txKeys = Object.keys(txs);

      const list: TransactionStatus[] = [];
      for (let key of txKeys) {
        const txStatus = txs[key];

        // Are only interested in pending txs with a symbol
        if (!txStatus.symbol || txStatus.state !== "accepted") continue;

        list.push(txStatus);
      }

      return list;
    });

    const txMatchesUnpegSymbol = (pegAssetSymbol: string) => (
      txStatus: TransactionStatus,
    ) => {
      return (
        txStatus.symbol?.toLowerCase() ===
        getUnpeggedSymbol(pegAssetSymbol.toLowerCase()).toLowerCase()
      );
    };

    const assetList = computed<TokenListItem[]>(() => {
      const balances =
        selectedTab.value === "External Tokens"
          ? store.wallet.eth.balances
          : store.wallet.sif.balances;

      const pegList = pendingPegTxList.value;

      let listedTokens = allTokens.value
        .filter(
          ({ symbol }) =>
            symbol
              .toLowerCase()
              .indexOf(searchText.value.toLowerCase().trim()) > -1,
        )
        .map((asset) => {
          const amount = balances.find(({ asset: { symbol } }) => {
            return asset.symbol.toLowerCase() === symbol.toLowerCase();
          });

          // Get pegTxs for asset
          const pegTxs = pegList
            ? pegList.filter(txMatchesUnpegSymbol(asset.symbol))
            : [];

          // Is the asset from a supported network
          const supported = getIsSupportedNetwork(asset);

          if (!amount) {
            return {
              amount: AssetAmount(asset, "0"),
              asset,
              pegTxs,
              supported,
            };
          }

          return {
            amount,
            asset,
            pegTxs,
            supported,
          };
        });

      const listedTokensSorted = sortAssetAmount(listedTokens);

      return listedTokensSorted;
    });

    // TODO: add to utils
    function shortenHash(hash: string) {
      const start = hash.slice(0, 7);
      const end = hash.slice(-7);
      return `${start}...${end}`;
    }

    return {
      shortenHash,
      assetList,
      searchText,
      peggedSymbol(unpeggedSymbol: string) {
        if (unpeggedSymbol.toLowerCase() === "erowan") {
          return "rowan";
        }
        return "c" + unpeggedSymbol;
      },

      unpeggedSymbol(peggedSymbol: string) {
        if (peggedSymbol.toLowerCase() === "rowan") {
          return "erowan";
        }
        return peggedSymbol.replace(/^c/, "");
      },

      onTabSelected({ selectedTitle }: { selectedTitle: string }) {
        selectedTab.value = selectedTitle;
      },
    };
  },
});
</script>

<template>
  <Layout>
    <div class="search-text">
      <SifInput
        gold
        placeholder="Search name or paste address"
        type="text"
        v-model="searchText"
      />
    </div>
    <Tabs :defaultIndex="1" @tabselected="onTabSelected">
      <Tab title="External Tokens" slug="external-tab">
        <AssetList :items="assetList" v-slot="{ asset }">
          <SifButton
            :disabled="!asset.supported"
            :to="`/peg/${asset.asset.symbol}/${peggedSymbol(
              asset.asset.symbol,
            )}`"
            primary
            :data-handle="'peg-' + asset.asset.symbol"
            >Peg</SifButton
          >
          <Tooltip v-if="!asset.supported" message="Network not supported">
            &nbsp;<Icon icon="info-box-black" />
          </Tooltip>
        </AssetList>
      </Tab>
      <Tab title="Sifchain Native" slug="native-tab">
        <AssetList :items="assetList">
          <template #default="{ asset }">
            <SifButton
              :to="`/peg/reverse/${asset.asset.symbol}/${unpeggedSymbol(
                asset.asset.symbol,
              )}`"
              primary
              :data-handle="'unpeg-' + asset.asset.symbol"
              >Unpeg</SifButton
            >
          </template>
          <template #annotation="{ pegTxs }">
            <span v-if="pegTxs.length > 0">
              <Tooltip>
                <template #message>
                  <p>You have the following pending transactions:</p>
                  <br />
                  <p v-for="tx in pegTxs" :key="tx.hash">
                    <a
                      :href="`https://etherscan.io/tx/${tx.hash}`"
                      :title="tx.hash"
                      target="_blank"
                      >{{ shortenHash(tx.hash) }}</a
                    >
                  </p></template
                >
                <template #default
                  >&nbsp;<span data-handle="pending-tx-marker" class="footnote"
                    >*</span
                  ></template
                >
              </Tooltip>
            </span>
          </template>
        </AssetList>
      </Tab>
    </Tabs>
    <ActionsPanel connectType="connectToAll" />
  </Layout>
</template>
<style lang="scss" scoped>
.search-text {
  margin-bottom: 1rem;
}
.footnote {
  font-family: Arial, Helvetica, sans-serif;
  font-weight: bold;
  font-style: normal;
  color: $c_gold_dark;
}
</style>
