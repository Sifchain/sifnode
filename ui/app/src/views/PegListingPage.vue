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
      <Tab title="Standard">
        <AssetList :items="assetList" v-slot="{ asset }">
          <SifButton @click="handlePegClicked(asset)" primary>Peg</SifButton>
        </AssetList>
      </Tab>
      <Tab title="Pegged">
        <AssetList :items="assetList" v-slot="{ asset }">
          <SifButton @click="handleUnpegClicked(asset)" primary
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
<script>
import Tab from "@/components/shared/Tab.vue";
import Tabs from "@/components/shared/Tabs.vue";
import Layout from "@/components/layout/Layout.vue";
import AssetList from "@/components/shared/AssetList.vue";
import SifInput from "@/components/shared/SifInput.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import SifButton from "@/components/shared/SifButton.vue";
import { useTokenListing } from "@/components/tokenSelector/useSelectToken";

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
  setup() {
    const { store, actions } = useCore();

    const searchText = ref("");
    const selectedTab = ref("Standard");

    const allTokens = computed(() => {
      if (selectedTab.value === "Standard") {
        return actions.peg.getEthTokens();
      }

      if (selectedTab.value === "Pegged") {
        return actions.peg.getSifTokens();
      }
    });

    const assetList = computed(() => {
      const balances =
        selectedTab.value === "Standard"
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

          if (!amount) return { amount: "", asset };

          return {
            amount: amount.toFixed(amount.asset.decimals === 0 ? 0 : 6),
            asset,
          };
        });
    });

    return {
      assetList,
      searchText,
      handlePegClicked(asset) {
        alert("Launch peg dialog");
      },
      handleUnpegClicked(asset) {
        alert("Launch unpeg dialog");
      },
      onTabSelected({ selectedTitle }) {
        selectedTab.value = selectedTitle;
      },
    };
  },
});
</script>