<template>
  <div class="df fdc" style="background: linear-gradient(90deg, black, transparent); padding: 16px">
    <div class="wallet-container mb8 df fdc aifs">

      <div class="df fdc w100" v-if="!sifWallet.isConnected">
        <textarea
          v-if="!sifWallet.isConnected"
          class="mb8"
          v-model="localMnemonic"
          placeholder="Mnemonic..."
        ></textarea>
        <div class="df">
          <button @click="submit" class="mr8">
            Connect Wallet
          </button>
          <button @click="reset">
            Clear
          </button>
        </div>
      </div>

      <div v-else class="df fdc aifs w100">
        <div class="df fdr address-container w100 mb8">
          <div class="df connected-dot mr8"></div>
          <div class="df fdc aifs">
            Address: {{sifWallet.address}}
          <div v-for="coin in sifWallet.balances" :key="coin.denom">
            Balance: {{coin.amount}}{{coin.denom}}
          </div>
        </div>
      </div>

        
        <button @click="reset">
          Clear
        </button>
      </div>


      <div style="color:salmon; font-weight: bold">{{errorMessage}}</div>

    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, reactive, readonly } from "vue";

import CSignIn from "@/components/CSignIn.vue"
// import { CWalletStore, WalletStore } from "../../../core/src/store/wallet"
import { signInCosmosWallet, getCosmosBalanceAction } from "../../../core/src/actions/CWalletActions"

export default defineComponent({
  name: "SifWallet",
  setup() {
    // local reactive variables
    const localMnemonic = ref("race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow")
    let errorMessage = ref()
    const initSifWallet = {
      isConnected: false,
      client: undefined,
      address: undefined,
      balances: undefined
    }
    const sifWallet = reactive({...initSifWallet})

  const activeSifAddress = readonly(sifWallet.client)
    // submit to actions
    async function submit() {
      errorMessage.value = ""
      if (!localMnemonic.value) { return errorMessage.value = "Mnemonic required to send" }

      try {
        // assign local store
        sifWallet.client = await signInCosmosWallet(localMnemonic.value.trim())
        sifWallet.address = sifWallet.client.senderAddress
        sifWallet.isConnected = true
        sifWallet.balances = (await getCosmosBalanceAction(sifWallet.address)).balance
      } catch(error) { 
        errorMessage.value = error 
      }

    }
    // reset
    async function reset() {
      localMnemonic.value = ""
      errorMessage.value = ""
      Object.assign(sifWallet, initSifWallet)
    }

    return {
      sifWallet,
      submit, 
      errorMessage,
      localMnemonic,
      reset
    }
  },
});
</script>

<style scoped>

.df {display: flex}
.fdc {flex-direction: column}
.aic {align-items: center;}
.aifs {align-items: flex-start;}

.w100 {width: 100%}
.mr4 {margin-right: 4px}
.mr8 {margin-right: 8px}
.mr12 {margin-right: 12px}
.mb4 {margin-bottom: 4px}
.mb8 {margin-bottom: 8px}
.mb12 {margin-bottom: 12px}

.pb4 {padding-bottom: 4px}
.pb8 {padding-bottom: 8px}
.pb12 {padding-bottom: 12px}


.wallet-container {
  width: 500px;
  font-family: monospace
}
.address-container { background: #363636; color: beige; padding: 6px; height: 100px}
.connected-dot { width: 10px; height: 10px; border-radius: 10px; margin-top: 4px; background: darkseagreen}
button { padding: 6px; width: fit-content; }
textarea {height: 100px; padding: 6px}
</style>