<template>
  <div class="sif-select">
    <div class="sif-select__activator" @click="openModal">
      <slot name="activator"></slot>
    </div>
    <div class="sif-select__wrapper" v-if="isOpen">
      <div class="sif-select__content">
        <div class="sif-select__close" @click="closeModal">&times;</div>
        <div class="sif-select__header">
          <h3 class="sif-select__title">Select a token</h3>
          <SifInput gold placeholder="Search name or paste address" />
        </div>
        <div class="sif-select__body">
          <slot></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import SifInput from '../Elements/SifInput.vue';

export default {
  components: {
    SifInput,
  },
  
  data() {
    return {
      isOpen: false,
    }
  },

  methods: {
    openModal() {
      this.isOpen = true;
    },
    closeModal() {
      this.isOpen = false;
    }
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

  &__content {
    position: relative;
    width: 410px;
    min-height: 60vh;
    max-height: 80vh;
    background: $c_white;
    border-radius: $br_md;
    box-shadow: $bs_default;
    padding-top: 30px;
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
}
</style>