<template>
  <div>
    <div id="nav">
      <router-link to="/">Wallet</router-link>
    </div>
    <router-view />
  </div>
</template>

<script lang="ts">
import { defineComponent, provide } from "vue";
import {
  createStore,
  createApi,
  getWeb3,
  createActions,
  getFakeTokens as getSupportedTokens,
} from "../../core";

export default defineComponent({
  name: "App",
  setup() {
    const api = createApi({
      getWeb3,
      getSupportedTokens,
    });

    const store = createStore();
    const actions = createActions({ store, api });
    console.log({ actions });
    provide("api", api);
    provide("store", store);
    provide("actions", actions);
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
