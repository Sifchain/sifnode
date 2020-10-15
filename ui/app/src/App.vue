<template>
  <div>
    Wallet store access in App component: {{CWalletStore}}

    <div id="nav">
      <router-link to="/">Wallet</router-link> |
      <router-link to="/cosmos">CWallet</router-link> |
      <router-link to="/list">List</router-link>
    </div>
    <router-view />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { CWalletStore } from "../../core/src/store/wallet"

export default defineComponent({
  name: "App",
  setup() {
    return {CWalletStore}
  }
});
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
