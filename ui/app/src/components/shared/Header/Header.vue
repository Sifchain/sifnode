<script lang="ts">
import { defineComponent } from "vue";
import Menu from "../Menu/Menu.vue";
import RowanPrice from "../RowanPrice.vue";
import TVL from "../TVL.vue";

export default defineComponent({
  components: {
    Menu,
    RowanPrice,
    TVL,
  },
});
</script>

<template>
  <header class="header">
    <div class="secondary-header-container">
      Welcome to our BetaNet! Please be aware that while this has passed several
      public security audits, it is still in Beta. We encourage you to first
      understand how it works before transacting on it and exercise caution at
      all times.
      <a
        target="_blank"
        href="https://docs.sifchain.finance/resources/betanet-launch"
        >Here's a list of security measures taken</a
      >, but please note that all transactions conducted are at your own risk.
      This platform is not open to US users.
    </div>
    <div class="primary-header-container">
      <div>
        <Menu />
      </div>
      <div class="third-container">
        <Suspense>
          <div>
            <TVL />
          </div>
        </Suspense>
        <Suspense>
          <div>
            <RowanPrice />
          </div>
        </Suspense>
        <div>
          <slot name="right"></slot>
        </div>
      </div>
    </div>
  </header>
</template>

<style lang="scss" scoped>
.header {
  width: 100%;
  position: fixed;
  top: 0;
  left: 0;
  background: $c_white;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  .primary-header-container {
    background: #000;
  }
  .primary-header-container,
  .secondary-header-container {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding: 1rem 2rem;
    font-size: 12px;
  }
  .secondary-header-container {
    display: block;
    text-align: left;
    background: $c_gold;
    color: $c_white;
    padding: 8px 2rem;
    a {
      color: $c_white;
      text-decoration: underline;
    }
  }
  .third-container {
    display: flex;
    justify-content: center;
    > div {
      margin-left: 15px;
    }
  }
}
</style>
