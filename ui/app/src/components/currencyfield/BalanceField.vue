<template>
  <span v-if="available !== null">Balance: {{ available }}</span>
</template>

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

    const available = computed(
      () =>
        balances.value
          .find((bal) => bal.asset.symbol === props.symbol)
          ?.toFixed(2) ?? "0"
    );

    return { available };
  },
});
</script>