import { useCore } from "@/hooks/useCore";
import { computed } from "@vue/reactivity";

const shorten = (num: number) => (str: string) => str.slice(0, num) + "...";

export function useWalletButton({ addrLen = 5 }: { addrLen: number }) {
  const { store } = useCore();

  const connected = computed(
    () => store.wallet.eth.isConnected || store.wallet.sif.isConnected
  );

  const connectedText = computed(() =>
    [store.wallet.eth.address, store.wallet.sif.address]
      .filter(Boolean)
      .map(shorten(addrLen))
      .join(", ")
  );

  return { connected, connectedText };
}
