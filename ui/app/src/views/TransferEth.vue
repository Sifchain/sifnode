<template>
  <div class="list">
    <div v-if="walletConnected">
      <input v-model="accountAddresText" placeholder="0x12347..." />
      <button @click="transfer">Transfer</button>
    </div>
    <div v-else>Your wallet is not connected</div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";

import { computed, ref } from "@vue/reactivity";
// import { getFakeTokens } from "../../../core";
import JSBI from "jsbi";
import { useCore } from "../hooks/useCore";

export default defineComponent({
  name: "ListPage",
  setup() {
    const { api, actions, store } = useCore();
    const accountAddresText = ref("");
    const walletConnected = computed(() => store.wallet.etheriumIsConnected);

    onMounted(async () => {
      await actions.init();
    });

    async function transfer() {
      if (accountAddresText.value === "")
        throw new Error("Account must be supplied");
      // const toks = await getFakeTokens();
      // const ATK = toks.find((tok) => tok.symbol === "ATK");
      // if (!ATK) throw new Error("ATK not found");
      const hash = await api.EtheriumService.transfer({
        amount: JSBI.BigInt("10000000000000000000"),
        recipient: accountAddresText.value,
        // asset: ATK,
      });
      console.log(hash);
    }

    return { transfer, walletConnected, accountAddresText };
  },
  components: {},
});
</script>
<style scoped>
table {
  margin: 0 auto;
}
</style>
