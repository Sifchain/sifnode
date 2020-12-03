<script lang="ts">
import { defineComponent } from "vue";

import Layout from "@/components/layout/Layout.vue";
export default defineComponent({
  components: { Layout },
});
</script>

<template>
  <Layout class="pool">
    <router-link class="button" to="/pool/create-pool">Create Pool</router-link>
    <router-link class="button" to="/pool/add-liquidity"
      >Add Liquidity</router-link
    >
    <router-link class="button" to="/pool/remove-liquidity"
      >Remove Liquidity</router-link
    >
  </Layout>
</template>

<style scoped>
.button {
  display: block;
  padding: 1rem;
  text-decoration: none;
  color: #333;
  background: #aaa;
  margin-bottom: 1rem;
}
.button:hover {
  color: #666;
  background: #bbb;
}
</style>