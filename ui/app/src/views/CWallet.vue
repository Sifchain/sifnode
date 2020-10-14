<template>
  <div class="home">
      <c-sign-in />
  </div>
</template>

<script lang="ts">
import { onMounted, ref } from 'vue'
import { api, entities } from '../../../core'
import CSignIn from '@/components/CSignIn'

export default {
  name: 'CWallet',
  components: {
    CSignIn
  },
  setup() {
  }
};
</script>
