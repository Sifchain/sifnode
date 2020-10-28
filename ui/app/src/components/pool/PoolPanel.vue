<template>
  <router-link to="/pool/create-pair">Create Pair</router-link>
  <router-link to="/pool/add-liquidity">Add Liquidity</router-link>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "PoolPanel",
});
</script>
