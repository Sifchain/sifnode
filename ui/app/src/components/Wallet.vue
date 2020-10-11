<template>
  <div class="home">
      {{JSON.stringify(balances)}}
  </div>
</template>

<script lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { api, entities } from '../../../core'

export default {
  name: 'Wallet',
  setup() {
    let balances = ref<entities.AssetAmount[]>([])
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
