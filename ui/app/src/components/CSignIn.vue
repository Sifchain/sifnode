<template>
  <div>
    <div class="container">
      <div class="row">
        <div class="button" v-if="address">
          <div class="button__text button__address">
            {{ address }}
          </div>
        </div>
        <div class="button" @click="buttonClick">
          <div class="button__icon">
            <c-icon-lock />
          </div>
          <div class="button__text">
            {{ address ? "Log out" : "Sign in " }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="container-dropdown">
          <div v-if="dropdown && !address">
            <div class="dropdown">
              <div class="dropdown__textarea">
                <textarea
                  v-model="mnemonic"
                  placeholder="Mnemonic..."
                  class="dropdown__textarea__input"
                ></textarea>
                <div
                  class="dropdown__textarea__icon"
                  @click="mnemonicGenerate()"
                >
                  <c-icon-magic />
                </div>
              </div>
              <div
                :class="[
                  'dropdown__button',
                  `button__disabled__${!mnemonicIsValid}`,
                ]"
                @click="mnemonicImport"
              >
                <div class="dropdown__button__text">
                  Import mnemonic
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import CIconLock from "./CIconLock";
import CIconMagic from "./CIconMagic";

import * as bip39 from "bip39";
export default {
  components: {
    CIconLock,
    CIconMagic,
  },
  props: ["foo"],
  data() {
    return {
      mnemonic: "",
      dropdown: false,
    };
  },
  computed: {
    mnemonicClean() {
      return this.mnemonic.trim();
    },
    mnemonicIsValid() {
      return bip39.validateMnemonic(this.mnemonicClean);
    },
    address() {
      const { client } = this.$store.state.cosmos;
      return client && client.senderAddress;
    },
  },
  methods: {
    buttonClick() {
      if (this.address) {
        this.$store.dispatch("cosmos/accountSignOut");
      } else {
        this.mnemonic = "";
        this.dropdown = !this.dropdown;
      }
    },
    async mnemonicImport() {
      if (this.mnemonicIsValid) {
        const mnemonic = this.mnemonicClean;
        await this.$store.dispatch("cosmos/accountSignIn", { mnemonic });
      }
    },
    mnemonicGenerate() {
      const mnemonic = bip39.generateMnemonic();
      this.mnemonic = mnemonic;
    },
    truncate(string) {
      return `${string.substring(0, 16)}...`;
    },
  },
};
</script>