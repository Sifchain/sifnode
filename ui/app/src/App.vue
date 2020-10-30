<template>
  <div class="main">
    <Header>
      <template v-slot:right>
        <WithWallet>
          <template v-slot:disconnected="{ requestDialog }">
            <SifButton primary @click="requestDialog">Connect Wallet</SifButton>
          </template>
          <template v-slot:connected="{ connectedText, requestDialog }">
            <SifButton primary @click="requestDialog">
              {{ connectedText }}
            </SifButton>
          </template>
        </WithWallet>
      </template>
    </Header>

    <router-view />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { useCore } from "./hooks/useCore";
import WithWallet from "@/components/wallet/WithWallet.vue";
import Header from "./components/shared/Header.vue";
import Main from "./components/shared/Main.vue";
import SifButton from "./components/shared/SifButton.vue";

export default defineComponent({
  name: "App",
  components: {
    Header,
    Main,
    WithWallet,
    SifButton,
  },

  setup() {
    const { actions } = useCore();
    onMounted(async () => {
      await actions.token.refreshTokens();
    });
  },
});
</script>

<style lang="scss">
#app,
#portal-target {
  font: italic normal bold 14px/22px $f_default;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
}

a {
  font-weight: bold;
}

.main {
  min-height: 100vh;
}
</style>