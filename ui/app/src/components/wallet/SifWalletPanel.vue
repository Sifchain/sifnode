<script lang="ts">
import { computed, defineComponent } from "vue";
import { ref, Ref } from "@vue/reactivity"; /* eslint-disable-line */
import { useCore } from "@/hooks/useCore";
import BalanceTable from "./BalanceTable.vue";
import SifButton from "@/components/shared/SifButton.vue";

function useCosmosWallet({
  error,
  mnemonic,
}: {
  error: Ref<string>;
  mnemonic: Ref<string>;
}) {
  const { store, actions } = useCore();

  async function handleDisconnectClicked() {
    try {
      await actions.wallet.disconnect();
      mnemonic.value = "";
    } catch (err) {
      error.value = err;
    }
  }

  async function handleConnectClicked() {
    try {
      await actions.wallet.connect(mnemonic.value);
      mnemonic.value = "";
    } catch (err) {
      error.value = err;
    }
  }

  const address = computed(() => store.wallet.sif.address);
  const balances = computed(() => store.wallet.sif.balances);
  const connected = computed(() => store.wallet.sif.isConnected);

  return {
    address,
    balances,
    connected,
    handleConnectClicked,
    handleDisconnectClicked,
  };
}

export default defineComponent({
  name: "SifWalletController",
  components: { BalanceTable, SifButton },
  setup() {
    const error = ref("");
    // TODO: remove hard coded mnemonic
    const mnemonic = ref("");

    const {
      balances,
      address,
      connected,
      handleDisconnectClicked,
      handleConnectClicked,
    } = useCosmosWallet({
      error,
      mnemonic,
    });

    const connectionStarted = ref<boolean>(false);

    function handleStartConnectClicked() {
      connectionStarted.value = true;
    }

    return {
      error,
      mnemonic,
      balances,
      address,
      connected,
      connectionStarted,
      handleDisconnectClicked: () => {
        handleDisconnectClicked();
        connectionStarted.value = false;
      },
      handleStartConnectClicked,
      handleConnectClicked,
      isLocalChain: !process.env.VUE_APP_SIFNODE_API,
    };
  },
});
</script>

<template>
  <div class="wrapper">
    <div v-if="connected">
      <p>{{ address }}</p>
      <BalanceTable :balances="balances" />
      <SifButton secondary @click="handleDisconnectClicked">
        Disconnect SifWallet
      </SifButton>
    </div>
    <div v-else>
      <div v-if="!connectionStarted">
        <SifButton secondary @click="handleStartConnectClicked"
          >Connect to SifChain</SifButton
        >
      </div>
      <div v-else>
        <div v-if="isLocalChain">
          <button
            @click="
              mnemonic =
                'race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow'
            "
          >
            Shadowfiend
          </button>
          <button
            @click="
              mnemonic =
                'hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard'
            "
          >
            Akasha
          </button>
        </div>
        <textarea
          class="textarea"
          v-model="mnemonic"
          placeholder="Mnemonic..."
        ></textarea
        ><br />
        <SifButton secondary @click="connectionStarted = false"
          >Clear</SifButton
        >
        <SifButton secondary @click="handleConnectClicked">Login</SifButton>
        <div v-if="error">{{ error }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.textarea {
  width: 100%;
  min-height: 50px;
}
</style>