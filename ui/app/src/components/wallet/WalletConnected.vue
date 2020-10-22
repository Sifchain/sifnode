<template>
  <p>{{ address }}</p>
  <table>
    <tr v-for="assetAmount in balances" :key="assetAmount.asset.symbol">
      <td align="left">{{ assetAmount.asset.symbol }}</td>
      <td align="right">{{ assetAmount.toFixed(2) }}</td>
    </tr>
  </table>
  <button @click="handleDisconnectClicked">DisconnectWallet</button>
  <button @click="handleClose">Close</button>
</template>

<script lang="ts">
import { computed, defineComponent } from "vue";
import { ref } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";

export default defineComponent({
  name: "WalletConnectSelector",
  setup(_, { emit }) {
    const { store, actions } = useCore();
    const connecting = ref<boolean>(false);

    async function handleDisconnectClicked() {
      await actions.disconnectWallet();
      emit("close");
    }

    function handleClose() {
      emit("close");
    }

    const balances = computed(() => store.wallet.eth.balances);
    const address = computed(() => store.wallet.eth.address);

    return {
      balances,
      address,
      handleDisconnectClicked,
      handleClose,
      connecting,
    };
  },
});
</script>