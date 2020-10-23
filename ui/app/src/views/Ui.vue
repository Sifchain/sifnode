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
          <div class="d-flex j-between">
            <IconButton label="Swap" active />
            <IconButton label="Swap" />
            <IconButton label="Swap" disabled />
            <IconButton label="Swap" disabled />
            <IconButton label="Swap" disabled />
          </div>
        </template>
        <div>
          <InputGroup>
            <CurrencyForm label="From" />
          </InputGroup>

          <InputGroup>
            <CurrencyForm label="To" />
          </InputGroup>

          <SwapPanel />

          <p class="text--small">No wallet connected</p>

          <SifButton primary block>Connect to Wallet</SifButton>
        </div>
      </SifApp>

  </div>
</template>

<script>
import SifButton from '../components/Elements/SifButton.vue';
import IconButton from '../components/Elements/IconButton.vue';
import InputGroup from '../components/Elements/InputGroup.vue';
import SifApp from '../components/Blocks/SifApp.vue';
import CurrencyForm from '../components/Blocks/CurrencyForm.vue';
import SwapPanel from '../components/Blocks/SwapPanel.vue';

export default {
  components: {
    SifButton,
    IconButton,
    SifApp,
    InputGroup,
    CurrencyForm,
    SwapPanel
  }
}
</script>

<style>

</style>