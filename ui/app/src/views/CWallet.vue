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

    <button @click="()=> {
      errorMessage = null
      localMnemonic = null
    }">
      clear
    </button>
      Local mnemonic value: {{localMnemonic}} <br>
      Core store: {{CWalletStore}}
    <div style="color:salmon; font-weight: bold">{{errorMessage}}</div>

  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed } from "vue";

import CSignIn from "@/components/CSignIn.vue"
import { CWalletStore, WalletStore } from "../../../core/src/store/wallet"
import { signInCosmosWallet } from "../../../core/src/actions/CWalletActions"

export default defineComponent({
  name: "CWallet",
  setup() {
    // local reactive variables
    const localMnemonic = ref()
    const errorMessage = ref()
    // submit to actions
    async function submit() {
      errorMessage.value = ""
      if (!localMnemonic.value) { 
        return errorMessage.value = "Mnemonic required to send" 
      }

      const parsedMnemonic = computed(() => {
        return (localMnemonic.value).trim()
      })

      try {
        await signInCosmosWallet(parsedMnemonic.value)
      } catch(error) { 
        errorMessage.value = error 
      }

    }

    return {
      CWalletStore, 
      submit, 
      errorMessage,
      localMnemonic
    }
  },
});
</script>

<style scoped>
.df {display: flex}
.fdc {flex-direction: column}
</style>