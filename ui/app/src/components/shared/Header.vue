<template>
  <header class="header">
    <div class="secondary-header-container">Welcome to our BetaNet! Please be aware that while this has passed several public security audits, it is still in Beta. We encourage you to first understand how it works before transacting on it and exercise caution at all times. <a target="_blank" href="https://docs.sifchain.finance/resources/betanet-launch">Here's a list of security measures taken</a>, but please note that all transactions conducted are at your own risk.</div>
    <div class="primary-header-container">
      <div>
        <Menu />
        <!-- <SifSwitch label="Typeface Cursive" v-model="typefaceCursive" /> -->
      </div>
      <div>
        <slot name="right"></slot>
      </div>
    </div>
  </header>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Menu from "./Menu.vue";
// import SifSwitch from "./SifSwitch.vue";

export default defineComponent({
  components: {
    Menu,
    // SifSwitch
  },

  data() {
    return {
      typefaceCursive: true,
    };
  },
});
</script>

<style lang="scss" scoped>
.header {
  width: 100%;
  position: fixed;
  top: 0;
  left: 0;
  background: $c_white;
  border-bottom: $divider;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  .primary-header-container, .secondary-header-container {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding: 1rem 2rem
  }
  .secondary-header-container {
    display: block;
    text-align:  center;
    background: $c_gold;
    color: $c_white;
    padding: 8px 2rem;
    a {
      color: $c_white;
      text-decoration: underline;
    }
  }
}
</style>