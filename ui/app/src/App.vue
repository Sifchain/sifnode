<template>
  <div class="main">
    <Header>
      <template v-slot:right>
        <WithWallet>
          <template v-slot:disconnected="{ requestDialog }">
            <Pill
              data-handle="button-connected"
              color="danger"
              @click="requestDialog"
            >
              Not connected
            </Pill>
          </template>
          <template v-slot:connected="{ requestDialog }">
            <Pill
              data-handle="button-connected"
              @click="requestDialog"
              color="success"
              class="connected-button"
              >CONNECTED</Pill
            >
          </template>
        </WithWallet>
      </template>
    </Header>
    <router-view />
    <Notifications />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import WithWallet from "@/components/wallet/WithWallet.vue";
import Header from "./components/shared/Header/Header.vue";
import Pill from "./components/shared/Pill/Pill.vue";
import Footer from "./components/shared/Footer/Footer.vue";
import SifButton from "./components/shared/SifButton.vue";
import Notifications from "./components/Notifications.vue";
import { useInitialize } from "./hooks/useInitialize";
export default defineComponent({
  name: "App",
  components: {
    Header,
    Notifications,
    WithWallet,
    SifButton,
    Footer,
    Pill,
  },
  setup() {
    /// Initialize app
    useInitialize();
  },
});
</script>

<style lang="scss">
@import "normalize-scss";

// Import fonts:
@import url("https://fonts.googleapis.com/css2?family=PT+Serif:ital,wght@0,400;0,700;1,400;1,700&family=Rouge+Script&display=swap");
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap");

@import "@/scss/utilities.scss";
@import "@/scss/reset.scss";

#app,
#portal-target,
#tooltip-target {
  font: normal bold 14px/22px $f_default;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
}

input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

/* Firefox */
input[type="number"] {
  -moz-appearance: textfield;
}

a {
  font-weight: bold;
}

.main {
  min-height: 100vh;
}

.connected-button {
  cursor: pointer;
}
</style>
