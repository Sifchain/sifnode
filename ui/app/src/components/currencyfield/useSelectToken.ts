import { ModalBus } from "@/components/modal/ModalBus";

import SelectTokenDialog from "./SelectTokenDialog.vue";

export function useSelectTokens() {
  async function handleClicked(localBalance: { symbol: string }) {
    ModalBus.emit("open", {
      component: SelectTokenDialog,
      props: { localBalance },
    });
  }

  return { handleClicked };
}
