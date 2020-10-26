<template>
  <div class="sif-select">
    <div class="sif-select__activator" @click="openMenu">
      <SifButton primary block v-if="!selected">Select</SifButton>
      <SifButton secondary block v-else>{{ selected.name }}</SifButton>
    </div>
    <div class="sif-select__wrapper" v-if="isOpen" @click="closeMenu">
      <div class="sif-select__content" @click.stop>
        <div class="sif-select__close" @click="closeMenu">&times;</div>
        <div class="sif-select__header">
          <h3 class="sif-select__title">Select a token</h3>
          <SifInput gold placeholder="Search name or paste address" />
          <h4 class="sif-select__list-title">Token Name</h4>
        </div>
        <div class="sif-select__body">
          <div 
            class="sif-select__option" 
            v-for="(token, index) in tokens" 
            :key="index"
            @click="selectOption($event, token)">
            <span>{{ token.name }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import SifInput from '../Elements/SifInput.vue';
import SifButton from '../Elements/SifButton.vue';

export default {
  components: {
    SifInput,
    SifButton,
  },

  props: {
    tokens: Array,
  },
  
  data() {
    return {
      isOpen: false,
      selected: null,
    }
  },

  methods: {
    openMenu() {
      this.isOpen = true;
    },
    closeMenu() {
      this.isOpen = false;
    },
    selectOption(event, token) {
      this.selected = token;
      this.$emit('change', event, token);
      this.closeMenu();
    },
  }
  
}
</script>

<style lang="scss">
.sif-select {
  &__wrapper {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: $zi_modal;
  }

  &__activator {
    width: 128px;
  }

  &__content {
    position: relative;
    width: 410px;
    min-height: 60vh;
    max-height: 80vh;
    padding-top: 30px;
    display: flex;
    flex-direction: column;
    background: $c_white;
    border-radius: $br_md;
    box-shadow: $bs_default;
  }

  &__close {
    position: absolute;
    top: 16px;
    right: 20px;
    font-size: 32px;
    font-weight: normal;
    color: $c_text;
    cursor: pointer;
  }

  &__header {
    padding: 16px;
  }

  &__title {
    font-size: $fs_lg;
    color: $c_text;
    margin-bottom: 1em;
    text-align: left;
  }

  &__list-title {
    color: $c_text;
    text-align: left;
    margin-top: 30px;
  }

  &__body {
    padding-top: 14px;
    flex-grow: 1;
    overflow-y: scroll;
    border-top: $divider;
    border-right: $divider;
  }

  &__option {
    margin-bottom: 22px;
    font-size: $fs_md;
    font-weight: bold;
    text-align: left;
    color: $c_text;
    padding-left: 15px;
    cursor: pointer;
  }
}
</style>