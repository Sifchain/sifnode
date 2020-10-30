<template>
  <div>
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

          <PriceCalculation />

          <p class="text--small">No wallet connected</p>

          <SifButton primary block>Connect to Wallet</SifButton>
        </div>
      </SifApp>

  </div>
</template>

<script>
import SifButton from '../components/shared/SifButton.vue';
import IconButton from '../components/shared/IconButton.vue';
import InputGroup from '../components/shared/InputGroup.vue';
import SifApp from '../components/shared/SifApp.vue';
import CurrencyForm from '../components/shared/CurrencyForm.vue';
import PriceCalculation from '../components/shared/PriceCalculation.vue';

export default {
  components: {
    SifButton,
    IconButton,
    SifApp,
    InputGroup,
    CurrencyForm,
    PriceCalculation
  }
}
</script>

<style>

</style>