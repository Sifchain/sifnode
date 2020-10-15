<template>
  <div class="df fdc">
    <textarea
      v-if="!CWalletStore.isConnected"
      v-model="localMnemonic"
      placeholder="Mnemonic..."
    ></textarea>
    <button @click="submit">
      Connect Wallet
    </button>
    {{localMnemonic}}
    {{CWalletStore}}
      <!-- <c-sign-in /> -->
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, isRef, computed,  toRef, isReactive } from "vue";

import CSignIn from "@/components/CSignIn.vue"
import { CWalletStore, WalletStore } from "../../../core/src/store/wallet"
import { connectToWallet } from "../../../core/src/actions/CWalletActions"

export default defineComponent({
  name: "CWallet",
  components: {
    CSignIn
  },
  setup() {
    const localMnemonic = ref()
    async function submit() {
      if (!localMnemonic) return
      const parsedMnemonic = computed(() => {
        // return (localMnemonic.value + "TRANSFORM").trim()
        return (localMnemonic.value).trim()

      })
      await connectToWallet(parsedMnemonic.value)
    }
    return {CWalletStore, submit, localMnemonic}
  },
});
</script>

<style scoped>
.df {display: flex}
.fdc {flex-direction: column}
</style>