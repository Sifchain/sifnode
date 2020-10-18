<template>
  <div class="list">
    <div v-if="walletConnected">
      <span>Transfer </span><input v-model="amount" type="number" /><span>
        ETH to
      </span>
      <input v-model="accountAddresText" placeholder="0x12347..." />
      <button @click="transfer">Transfer</button>
    </div>
    <div v-else>Your wallet is not connected</div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { computed, ref } from "@vue/reactivity";
import JSBI from "jsbi";
import { useCore } from "../hooks/useCore";

export default defineComponent({
  name: "ListPage",
  setup() {
    const { api, store } = useCore();
    const accountAddresText = ref("");
    const walletConnected = computed(() => store.wallet.etheriumIsConnected);
    const amount = ref(10);

    async function transfer() {
      if (accountAddresText.value === "")
        throw new Error("Account must be supplied");

      const hash = await api.EtheriumService.transfer({
        amount: JSBI.multiply(
          JSBI.BigInt(amount.value),
          JSBI.BigInt("1000000000000000000")
        ),
        recipient: accountAddresText.value,
      });
      console.log(hash);
    }

    return { transfer, amount, walletConnected, accountAddresText };
  },
  components: {},
});
</script>
<style scoped>
table {
  margin: 0 auto;
}
</style>
