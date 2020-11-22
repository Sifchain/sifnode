<template>
  <div class="list">
    <div v-if="walletConnected">
      <div>
        <span>Transfer </span><input v-model="amount" type="number" /><span>
          ETH to
        </span>
        <input v-model="accountAddressText" placeholder="0x12347..." />
        <button @click="transfer">Transfer</button>
      </div>
      <div>
        <span>Transfer </span><input v-model="amountATK" type="number" /><span>
          ATK to
        </span>
        <input v-model="tokenAccountAddress" placeholder="0x12347..." />
        <button @click="transferATK">Transfer</button>
      </div>
    </div>
    <div v-else>Your wallet is not connected</div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { computed, ref } from "@vue/reactivity";
import { useCore } from "../hooks/useCore";

import B from "ui-core";
import { getFakeTokens } from "ui-core";

export default defineComponent({
  name: "ListPage",
  setup() {
    const { api, actions, store } = useCore();
    const accountAddressText = ref("");
    const tokenAccountAddress = ref("");
    const amountATK = ref(900);
    const walletConnected = computed(() => store.wallet.eth.isConnected);
    const amount = ref(10);

    // Utility function to get ATK token
    async function getATK() {
      const tokens = await getFakeTokens();
      const ATK = tokens.find(({ symbol }) => symbol === "ATK");
      if (!ATK) throw new Error("doesnt return ATK");
      return ATK;
    }

    async function transfer() {
      if (accountAddressText.value === "")
        throw new Error("Account must be supplied");

      const hash = await actions.ethWallet.transferEthWallet(
        amount.value,
        accountAddressText.value
      );

      console.log(hash);
    }

    async function transferATK() {
      console.log("transferATK");
      if (tokenAccountAddress.value === "")
        throw new Error("Account must be supplied");

      const ATK = await getATK();

      const hash = await actions.ethWallet.transferEthWallet(
        amountATK.value,
        tokenAccountAddress.value,
        ATK
      );

      console.log(hash);
    }

    return {
      transfer,
      transferATK,
      amount,
      amountATK,
      walletConnected,
      tokenAccountAddress,
      accountAddressText,
    };
  },
  components: {},
});
</script>
<style scoped>
table {
  margin: 0 auto;
}
</style>
