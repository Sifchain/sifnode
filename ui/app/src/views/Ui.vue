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
          <InputGroup>
            <CurrencyForm label="From" />
          </InputGroup>

          <InputGroup>
            <CurrencyForm label="To" />
          </InputGroup>

          <p class="text--small">No wallet connected</p>

          <SifButton primary block>Connect to Wallet</SifButton>
        </div>
      </SifApp>

  </div>
</template>

<script>
import SifButton from '../components/Elements/SifButton.vue';
import InputGroup from '../components/Elements/InputGroup.vue';
import SifApp from '../components/Blocks/SifApp.vue';
import CurrencyForm from '../components/Blocks/CurrencyForm.vue';

export default {
  components: {
    SifButton,
    SifApp,
    InputGroup,
    CurrencyForm,
  }
}
</script>

<style>

</style>