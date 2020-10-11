<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
// import {api} from '../../../core'
import * as x from '../../../core'

import HelloWorld from '@/components/HelloWorld.vue'; // @ is an alias to /src

export default defineComponent({
  name: 'Wallet',
  setup(props) {
    console.log(x)
  }
});
</script>
