import { ModalBus } from "@/components/modal/ModalBus";

import SelectTokenDialog from "./SelectTokenDialog.vue";

export function useSelectTokens() {
  async function handleClicked(label: string) {
    console.log("Sending key to dialog: " + label);
    ModalBus.emit("open", {
      component: SelectTokenDialog,
      props: { label },
    });
  }

  return { handleClicked };
}
