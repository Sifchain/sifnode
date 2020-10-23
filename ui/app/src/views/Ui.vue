<template>
  <div>
    <div class="mt-6">
      <SifButton primary>Primary</SifButton>
      <SifButton secondary>Secondary</SifButton>
      <SifButton disabled>Disabled</SifButton>
    </div>

    <div class="mt-6">
      <SifButton primary block>Primary Block</SifButton>
    </div>
  </div>
</template>

<script>
import SifButton from '../components/Elements/SifButton'

export default {
  components: {
    SifButton,
  }
}
</script>

<style>

</style>