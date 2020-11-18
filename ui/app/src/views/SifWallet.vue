<template>
  <div class="df fdc sifwallet-container">
    <div class="wallet-container mb8 df fdc aifs">
      <div
        class="df fdc w100 mb8"
        v-if="store.wallet.sif.balances.length === 0"
      >
        <!-- Best way here is to have input address, to get balances first,
             then if want to transact, add mnemonic. Also need to add create
        -->
        <p class="mb8">
          You may either input your pubkey to get balances. Or input your
          mnemonic to get balances and be ready to sign a transaction. You may
          add mnemonic later.
        </p>
        <input
          class="mb8 monospace"
          placeholder="sifAddress (todo)"
          v-model="store.wallet.sif.address"
        />

        <textarea
          v-if="!store.wallet.sif.isConnected"
          class="mb8"
          v-model="localMnemonic"
          placeholder="Mnemonic..."
        ></textarea>
        <div class="df">
          <button @click="signIn" class="mr8">Connect Wallet</button>
        </div>
      </div>

      <div v-else class="df fdc aifs w100">
        <div class="df fdr address-container w100 mb8">
          <div
            v-if="store.wallet.sif.isConnected"
            class="df connected-dot mr8"
          ></div>
          <div class="df fdc aifs">
            Address: {{ store.wallet.sif.address }}
            <div v-for="coin in store.wallet.sif.balances" :key="coin.denom">
              Balance: {{ coin.amount.toString() }} {{ coin.asset.symbol }}
            </div>
          </div>
        </div>

        <input
          class="mb8 monospace w100"
          placeholder="sendTo"
          v-model="sendTo"
        />
        <input
          class="mb8 monospace w100"
          placeholder="amount"
          v-model="amount"
        />

        <div class="mb8 w100" v-if="!store.wallet.sif.isConnected">
          <p class="mb8">Input mnemonic below to sign transaction.</p>
          <textarea
            class="mb8 w100"
            v-model="localMnemonic"
            placeholder="Mnemonic..."
          ></textarea>
          <div class="df mb8">
            <button @click="signIn" class="mr8">Connect Wallet</button>
          </div>
        </div>
        <button v-if="store.wallet.sif.isConnected" @click="send" class="mb8">
          Send
        </button>
      </div>
      <button @click="reset">Clear</button>
      <div style="color: salmon; font-weight: bold">{{ errorMessage }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { ref } from "@vue/reactivity";
import { useCore } from "../hooks/useCore";
import { Coin, Network } from "ui-core";
import JSBI from "jsbi";

export default defineComponent({
  name: "SifWallet",
  setup() {
    const { actions, store } = useCore();

    // local reactive variables
    const localMnemonic = ref(
      "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"
    );

    const errorMessage = ref();

    const amount = ref(50);
    const sendTo = ref("sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5");

    async function send() {
      if (!store.wallet.sif.isConnected) {
        return (errorMessage.value = "Not connected");
      }

      await actions.wallet.sendCosmosTransaction({
        asset: Coin({
          symbol: "catk",
          decimals: 18,
          name: "Nametoken",
          network: Network.SIFCHAIN,
        }),
        amount: JSBI.BigInt(amount.value),
        recipient: sendTo.value,
      });
    }

    async function signIn() {
      errorMessage.value = "";
      if (!localMnemonic.value) {
        return (errorMessage.value = "Mnemonic required to send");
      }
      try {
        await actions.wallet.connect(localMnemonic.value.trim());
      } catch (error) {
        errorMessage.value = error;
      }
    }

    async function reset() {
      localMnemonic.value = "";
      errorMessage.value = "";
      actions.wallet.disconnect();
    }

    return {
      store,
      sendTo,
      amount,
      signIn,
      errorMessage,
      localMnemonic,
      reset,
      send,
    };
  },
});
</script>

<style scoped>
.sifwallet-container {
  text-align: left;
}
.df {
  display: flex;
}
.fdc {
  flex-direction: column;
}
.aic {
  align-items: center;
}
.aifs {
  align-items: flex-start;
}

.w100 {
  width: 100%;
}
.mr4 {
  margin-right: 4px;
}
.mr8 {
  margin-right: 8px;
}
.mr12 {
  margin-right: 12px;
}
.mb4 {
  margin-bottom: 4px;
}
.mb8 {
  margin-bottom: 8px;
}
.mb12 {
  margin-bottom: 12px;
}

.pb4 {
  padding-bottom: 4px;
}
.pb8 {
  padding-bottom: 8px;
}
.pb12 {
  padding-bottom: 12px;
}

.monospace {
  font-family: monospace;
}
.tal {
  text-align: left;
}
.wallet-container {
  width: 500px;
  font-family: monospace;
}
.address-container {
  background: #363636;
  color: beige;
  padding: 6px;
  height: 100px;
}
.connected-dot {
  width: 10px;
  height: 10px;
  border-radius: 10px;
  margin-top: 4px;
  background: darkseagreen;
}
button {
  padding: 6px;
  width: fit-content;
}
textarea {
  height: 100px;
  padding: 6px;
}
p {
  margin: 0;
}
</style>