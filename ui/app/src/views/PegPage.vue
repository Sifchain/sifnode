<template>
  <Layout class="peg">
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
        <AssetList :items="assetList" />
      </Tab>
      <Tab title="Pegged">
        <AssetList :items="assetList" />
      </Tab>
    </Tabs>
    <ActionsPanel />
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

    // TODO: get balances and interleave balances
    const assetList = computed(() => {
      return allTokens.value
        .filter(
          ({ symbol }) =>
            symbol
              .toLowerCase()
              .indexOf(searchText.value.toLowerCase().trim()) > -1
        )
        .map((asset) => ({ amount: "", asset }));
    });

    return {
      assetList,
      searchText,
      handleNextStepClicked() {
        console.log("Next actions");
      },
      onTabSelected({ selectedTitle }) {
        selectedTab.value = selectedTitle;
      },
    };
  },
});
</script>