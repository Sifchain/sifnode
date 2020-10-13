<template>
  <div class="home">
    <Wallet/>
  </div>
</template>

<script lang="ts">

import Wallet from '@/components/Wallet.vue'

import { defineComponent, provide } from 'vue';
import { createStore,createApi, getWeb3, createUsecases, getFakeTokens } from '../../../core';

export default defineComponent({
  name: 'App',
  setup() { 
    const store = createStore()
    const { state } = store;
    const api = createApi({ getWeb3, getSupportedTokens: getFakeTokens })
    const usecases = createUsecases({ store, state, api})

    provide('api', api) 
    provide('state', store.state)
    provide('usecases', usecases)

  },
  components: {
    Wallet
  }
  
});
</script>
