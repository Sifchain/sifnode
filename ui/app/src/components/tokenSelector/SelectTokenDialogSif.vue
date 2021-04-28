<script lang="ts">
import { Asset } from "ui-core";
import { defineComponent, PropType } from "vue";
import { useCore } from "../../hooks/useCore";
import { generateTokenSearchLists } from "./tokenLists";
import SelectTokenDialog from "./SelectTokenDialog.vue";

export default defineComponent({
  components: { SelectTokenDialog },
  props: {
    selectedTokens: Array as PropType<string[]>,
    mode: { type: String, default: "from" },
  },
  emits: ["tokenselected"],
  setup(props, context) {
    const { config } = useCore();

    const { displayList, fullSearchList } = generateTokenSearchLists({
      walletLimit: 500,
      walletTokens: config.assets.filter((tok) => tok.network === "sifchain"),
    });

    function selectToken(symbol: string) {
      context.emit("tokenselected", symbol);
    }

    return { selectToken, displayList, fullSearchList };
  },
});
</script>
<template>
  <SelectTokenDialog
    :displayList="displayList"
    :fullSearchList="fullSearchList"
    :selectedTokens="selectedTokens"
    @tokenselected="selectToken"
  />
</template>
