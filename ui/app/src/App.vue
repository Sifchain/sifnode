<template>
  <div>
    <Header />
    <Main>
      <div>
        <p>Proof of technical capability</p>
        <div id="nav">
          <!-- <router-link to="/list">List</router-link> | -->
          <router-link to="/">EthWallet</router-link> |
          <router-link to="/ethtransfer">EthTransfer</router-link> |
          <router-link to="/sifwallet">SifWallet</router-link>
          <!-- <router-link to="/siftransfer">SifTransfer</router-link> -->
          <!-- <router-link to="/swap">Swap</router-link> -->
        </div>
        <router-view />
      </div>
    </Main>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { useCore } from "./hooks/useCore";
import Header from './components/Layout/Header.vue';
import Main from './components/Layout/Main.vue';

export default defineComponent({
  name: "App",
  components: {
    Header,
    Main,
  },

  setup() {
    const { actions } = useCore();
    onMounted(async () => {
      await actions.refreshTokens();
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

</style>
