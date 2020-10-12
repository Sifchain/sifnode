<template>
  <div class="home">
    <HelloWorld msg="Welcome"/>
    <Wallet/>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import HelloWorld from '@/components/HelloWorld.vue'; // @ is an alias to /src
import Wallet from '@/components/Wallet.vue';

export default defineComponent({
  name: 'Home',
  components: {
    HelloWorld,
    Wallet
  },
});
</script>
