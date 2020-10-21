<template>
  <div>
    {{message}}
  </div>
</template>

<script lang="ts">
// import CIconLock from "./CIconLock";
// import CIconMagic from "./CIconMagic";
// import * as bip39 from "bip39";
import { defineComponent } from "vue";

export default defineComponent({
  name: "CSignIn",
    data() {
      return {
        message: 'Hello!'
     }
  }
})
</script>