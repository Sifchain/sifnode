<script lang="ts">
import { useCore } from "@/hooks/useCore";
import { computed } from "@vue/reactivity";
import { defineComponent } from "vue";
export default defineComponent({
  props: ["symbol"],
  setup(props) {
    const { store } = useCore();

    const balances = computed(() => {
      return [...store.wallet.eth.balances, ...store.wallet.sif.balances];
    });

    const available = computed(() => {
      const found = balances.value.find(
        (bal) => bal.asset.symbol === props.symbol
      );
      if (!found) return "0";

      return found.toFormatted({
        decimals: Math.min(found.asset.decimals, 2),
        separator: true,
      });
    });

    return { available };
  },
});
</script>

<template>
  <span v-if="available !== '0'">Available: {{ available }}</span>
  <span v-else>&nbsp;</span>
</template>