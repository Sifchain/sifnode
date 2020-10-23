<template>
  <div>
    <div class="mt-6">
      <SifButton primary>Primary</SifButton>
      <SifButton secondary>Secondary</SifButton>
      <SifButton disabled>Disabled</SifButton>
    </div>

    <div class="mt-6 mb-6">
      <SifButton primary block>Primary Block</SifButton>
    </div>

    <!-- SifCard -->
      <SifApp>
        <template v-slot:header>
          HEADER
        </template>
        <div>
          <p class="text--small">No wallet connected</p>
          <SifButton primary block>Connect to Wallet</SifButton>
        </div>
      </SifApp>

  </div>
</template>

<script>
import SifButton from '../components/Elements/SifButton.vue';
import SifApp from '../components/Blocks/SifApp.vue';

export default {
  components: {
    SifButton,
    SifApp,
  }
}
</script>

<style>

</style>