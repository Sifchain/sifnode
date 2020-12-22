<template>
  <Layout class="peg">
    <Tabs>
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
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import { useTokenListing } from "@/components/tokenSelector/useSelectToken";

import { useCore } from "@/hooks/useCore";
import { defineComponent, ref } from "vue";
export default defineComponent({
  components: {
    Tab,
    Tabs,
    AssetList,
    Layout,
    ActionsPanel,
  },
  setup() {
    const { store } = useCore();
    const searchText = ref("");
    const { filteredTokens } = useTokenListing({
      searchText,
      store,
      tokenLimit: 20,
      walletLimit: 10,
    });

    console.log(filteredTokens.value);

    return {
      filteredTokens,
      searchText,
      handleNextStepClicked() {
        console.log("Next actions");
      },
    };
  },
});
</script>