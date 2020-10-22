<template>
  <Layout class="swap">
    <SwapPanel />
  </Layout>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import SwapPanel from "@/components/swap/SwapPanel.vue";
import Layout from "@/components/layout/Layout.vue";
export default defineComponent({
  name: "SwapPage",
  components: { Layout, SwapPanel },
});
</script>
