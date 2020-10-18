<template>
  <div>
    <div id="nav">
      <router-link to="/">Wallet</router-link> |
      <router-link to="/list">List all balances</router-link> |
      <router-link to="/transfer-ethereum"
        >Ethereum Network Transfer
      </router-link>
    </div>
    <router-view />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { useCore } from "./hooks/useCore";
export default defineComponent({
  name: "App",
  setup() {
    const { actions } = useCore();
    onMounted(async () => {
      await actions.refreshTokens();
    });
  },
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
