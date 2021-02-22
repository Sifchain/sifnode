<script>
import { defineComponent } from "vue";
import { RampInstantSDK } from "@ramp-network/ramp-instant-sdk";

export default defineComponent({
  methods: {
    openRamp: function() {
      new RampInstantSDK({
        swapAsset: "CELO,CUSD,DAI,DOT,ETH,KSM,USDC,USDT,XDAI,XLM",
        hostAppName: "Maker DAO",
        hostLogoUrl:
          "https://cdn-images-1.medium.com/max/2600/1*nqtMwugX7TtpcS-5c3lRjw.png"
      }).show();
    }
  }
});
</script>

<template>
  <div class="menu">
    <div>
      <router-link class="logo" to="/">S</router-link>
    </div>
    <div class="item">
      <a
        class="link"
        target="_blank"
        href="https://docs.sifchain.finance/resources/sifchain-dex-ui"
        >Discover</a
      >
    </div>
    <div class="item">
      <a
        class="link"
        target="_blank"
        href="https://docs.sifchain.finance/resources/faq#how-can-i-acquire-rowan-and-or-erowan"
        >Get Rowan</a
      >
    </div>

    <div class="item">
      <span class="link" v-on:click="openRamp">Buy Crypto</span>
    </div>

    <div class="item">
      <router-link class="link" to="/stats">Pool Stats</router-link>
    </div>
    <div class="item">
      <router-link class="link" to="/rewards">Stake & Delegate</router-link>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.menu {
  display: flex;
  align-items: baseline;
  .item {
    margin-right: 35px;
  }

  .logo {
    @include resetLink;
    font-family: $f_special;
    font-size: 40px;
    font-weight: normal;
    margin-right: 46px;
  }

  .link {
    @include resetLink;
    font-size: $fs_lg;
    font-weight: normal;
    opacity: 1;
    transition: opacity $trans_fast;
    cursor: pointer;
    &:hover {
      opacity: 0.5;
    }

    &.active {
      font-weight: bold;
      pointer-events: none;
    }
  }
}
</style>
