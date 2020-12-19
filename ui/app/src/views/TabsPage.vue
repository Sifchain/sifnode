<template>
  <Layout class="swap">
    <Tabs>
      <Tab title="Tab 1">Hello From Tab 1</Tab>
      <Tab title="Tab 2">Hello From Tab 2</Tab>
      <Tab title="Tab 3">Hello From Tab 3</Tab>
      <Tab title="Tab 4">Hello From Tab 4</Tab>
    </Tabs>
  </Layout>
</template>


<script>
import Tab from "@/components/shared/Tab.vue";
import Tabs from "@/components/shared/Tabs.vue";
import Layout from "@/components/layout/Layout.vue";

export default {
  components: {
    Tab,
    Tabs,
    Layout,
  },
};
</script>