import { ModalBus } from "@/components/modal/ModalBus";
import { useCore } from "@/hooks/useCore";
import { computed } from "@vue/reactivity";
import WalletConnectDialog from "./WalletConnectDialog.vue";

const shorten = (num: number) => (str: string) => str.slice(0, num) + "...";

export function useWalletButton({ addrLen = 5 }: { addrLen: number }) {
  const { store } = useCore();

  async function handleClicked() {
    ModalBus.emit("open", { component: WalletConnectDialog });
  }

  const connected = computed(
    () => store.wallet.eth.isConnected || store.wallet.sif.isConnected
  );

  const connectedText = computed(() =>
    [store.wallet.eth.address, store.wallet.sif.address]
      .filter(Boolean)
      .map(shorten(addrLen))
      .join(", ")
  );

  return { connected, connectedText, handleClicked };
}
