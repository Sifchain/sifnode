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
  createUsecases,
  getFakeTokens,
} from "../../core";

export default defineComponent({
  name: "App",
  setup() {
    const store = createStore();
    const { state } = store;
    const api = createApi({ getWeb3, getSupportedTokens: getFakeTokens });
    const usecases = createUsecases({ store, state, api });

    provide("api", api);
    provide("state", store.state);
    provide("usecases", usecases);
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
