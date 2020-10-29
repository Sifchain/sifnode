<template>
  <div>
    <div class="header-bar">
      <div class="header-nav">
        <router-link to="/">🏠</router-link>
        <router-link to="/">Discover</router-link>
        <router-link to="/">New wallet</router-link>
        <router-link to="/">Convert</router-link>
        <router-link to="/">Deposit</router-link>
      </div>
      <WithWallet>
        <template v-slot:disconnected="{ requestDialog }">
          <button @click="requestDialog">Connect Wallet</button>
        </template>
        <template v-slot:connected="{ connectedText, requestDialog }">
          <button @click="requestDialog">
            {{ connectedText }}
          </button>
        </template>
      </WithWallet>
    </div>
    <router-view />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { useCore } from "./hooks/useCore";
import WithWallet from "@/components/wallet/WithWallet.vue";
import Header from './components/shared/Header.vue';
import Main from './components/shared/Main.vue';

export default defineComponent({
  name: "App",
  components: {
    Header,
    Main,
    WithWallet
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
#app {
  font: italic normal bold 14px/22px $f_default;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
}

a {
  font-weight: bold;
}

.header-bar {
  height: 4rem; /* TODO: dynamize */
  padding: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  box-sizing: border-box;
  border-bottom: 1px solid #ccc;
}
.header-nav > * {
  margin-right: 1rem;
}
</style>