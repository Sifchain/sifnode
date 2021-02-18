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
    <Tabs @tabselected="onTabSelected">
      <Tab title="External Tokens">
        <AssetList :items="assetList" v-slot="{ asset }">
          <SifButton
            :to="`/peg/${asset.asset.symbol}/${peggedSymbol(
              asset.asset.symbol
            )}`"
            primary
            >Peg</SifButton
          >
        </AssetList>
      </Tab>
      <Tab title="Sifchain Native">
        <AssetList :items="assetList" v-slot="{ asset }">
          <SifButton
            :to="`/peg/reverse/${asset.asset.symbol}/${unpeggedSymbol(
              asset.asset.symbol
            )}`"
            primary
            >Unpeg</SifButton
          >
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
</style>
<script lang="ts">
import Tab from "@/components/shared/Tab.vue";
import Tabs from "@/components/shared/Tabs.vue";
import Layout from "@/components/layout/Layout.vue";
import AssetList from "@/components/shared/AssetList.vue";
import SifInput from "@/components/shared/SifInput.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import SifButton from "@/components/shared/SifButton.vue";

import { useCore } from "@/hooks/useCore";
import { defineComponent, ref } from "vue";
import { computed } from "@vue/reactivity";
export default defineComponent({
  components: {
    Tab,
    Tabs,
    AssetList,
    Layout,
    SifButton,
    SifInput,
    ActionsPanel,
  },
  setup(_, context) {
    const { store, actions } = useCore();

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

    const assetList = computed(() => {
      const balances =
        selectedTab.value === "External Tokens"
          ? store.wallet.eth.balances
          : store.wallet.sif.balances;

      return allTokens.value
        .filter(
          ({ symbol }) =>
            symbol
              .toLowerCase()
              .indexOf(searchText.value.toLowerCase().trim()) > -1
        )
        .map((asset) => {
          const amount = balances.find(({ asset: { symbol } }) => {
            return asset.symbol.toLowerCase() === symbol.toLowerCase();
          });

          if (!amount) return { amount: 0, asset };

          return {
            amount,
            asset,
          };
        });
    });

    return {
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
