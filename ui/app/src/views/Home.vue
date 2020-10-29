<template>
  <div class="home">
    <Wallet />
  </div>
</template>

<script lang="ts">
import Wallet from "@/components/Wallet.vue";

import { defineComponent } from "vue";

export default defineComponent({
  name: "HomePage",
  components: {
    Wallet,
  },
});
</script>
