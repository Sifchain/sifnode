<template>
  <div class="home">
      {{JSON.stringify(balances)}}
  </div>
</template>

<script lang="ts">
import { onMounted, ref } from 'vue'
import { api, entities } from '../../../core'

export default {
  name: 'CWallet',
  setup() {
    const balances = ref<entities.AssetAmount[]>([])
    const getAssetBalance = async () => {
      balances.value = await api.walletService.getAssetBalances()
    }
    // BigInt Error
    onMounted(getAssetBalance)
    return {
      balances
    }
  }
};
</script>
