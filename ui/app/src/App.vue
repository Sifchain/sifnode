<template>
  <div>
    <div class="header-bar">
      <div class="header-nav">
        <router-link to="/swap">Swap</router-link>
        <router-link to="/pool">Pool</router-link>
      </div>
      <WalletButton />
    </div>
    <router-view />
    <ModalRoot />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { useCore } from "./hooks/useCore";
import ModalRoot from "@/components/modal/ModalRoot.vue";
import WalletButton from "@/components/wallet/WalletButton.vue";

export default defineComponent({
  name: "App",
  setup() {
    const { actions } = useCore();
    onMounted(async () => {
      await actions.refreshTokens();
    });
  },
  components: { ModalRoot, WalletButton },
});
</script>

<style>
body {
  padding: 0;
  margin: 0;
}
</style>
<style scoped>
.header-bar {
  padding: 1rem;
  display: flex;
  justify-content: space-between;
  width: 100%;
  box-sizing: border-box;
  border-bottom: 1px solid #ccc;
}
.header-nav > * {
  margin-right: 1rem;
}
</style>