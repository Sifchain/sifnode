<script lang="ts">
import { defineComponent } from "vue";
import Layout from "@/components/layout/Layout.vue";
import { computed, ref, toRefs } from "@vue/reactivity";
import { useCore } from "@/hooks/useCore";
import { Asset, SwapState, useSwapCalculator } from "ui-core";
import { useWalletButton } from "@/components/wallet/useWalletButton";
import CurrencyPairPanel from "@/components/currencyPairPanel/Index.vue";
import Modal from "@/components/shared/Modal.vue";
import SelectTokenDialog from "@/components/tokenSelector/SelectTokenDialog.vue";
import PriceCalculation from "@/components/shared/PriceCalculation.vue";
import ActionsPanel from "@/components/actionsPanel/ActionsPanel.vue";
import ModalView from "@/components/shared/ModalView.vue";
import ConfirmationDialog, {
  ConfirmState,
} from "@/components/confirmationDialog/ConfirmationDialog.vue";
import { useCurrencyFieldState } from "@/hooks/useCurrencyFieldState";
import DetailsPanel from "@/components/shared/DetailsPanel.vue";
import { useRouter } from "vue-router";

export default defineComponent({
  components: {
    Layout,
  },

  setup(_, context) {
    const router = useRouter();
    const mode = computed(() => {
      return router.currentRoute.value.path.indexOf("unpeg") > -1
        ? "unpeg"
        : "peg";
    });
    return {
      mode,
    };
  },
});
</script>

<template>
  <Layout backLink="/peg">
    <div>{{ mode }} Asset</div>
  </Layout>
</template>