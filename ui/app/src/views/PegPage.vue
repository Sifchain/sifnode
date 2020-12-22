<template>
  <Layout class="peg">
    <SifInput
      gold
      placeholder="Search name or paste address"
      class="sif-input"
      type="text"
      v-model="searchText"
    />
    <Tabs @tabselected="onTabSelected">
      <Tab title="Standard">
        <AssetList :tokens="filteredTokens" />
      </Tab>
      <Tab title="Pegged">
        <AssetList :tokens="filteredTokens" />
      </Tab>
    </Tabs>
    <ActionsPanel />
  </Layout>
</template>

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

    const filteredTokens = computed(() => {
      return allTokens.value.filter(
        ({ symbol }) =>
          symbol.toLowerCase().indexOf(searchText.value.toLowerCase().trim()) >
          -1
      );
    });

    return {
      filteredTokens,
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