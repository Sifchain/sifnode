<template>
  <Layout class="swap">
    <SwapPanel />
  </Layout>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import SwapPanel from "@/components/swap/Index.vue";
import Layout from "@/components/layout/Layout.vue";
export default defineComponent({
  components: { Layout, SwapPanel },
});
</script>
