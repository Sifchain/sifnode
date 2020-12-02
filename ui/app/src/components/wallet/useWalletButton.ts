import { useCore } from "@/hooks/useCore";
import { computed } from "@vue/reactivity";

const shorten = (num: number) => (str: string) => str.slice(0, num) + "...";

export function useWalletButton(props?: { addrLen?: number }) {
  const { store } = useCore();

  const connected = computed(
    () => store.wallet.eth.isConnected || store.wallet.sif.isConnected
  );

  const connectedText = computed(() => {
    const addresses = [
      store.wallet.eth.address,
      store.wallet.sif.address,
    ].filter(Boolean);

    const addrLen = props?.addrLen || 10;
    return addresses
      .map(shorten(Math.round(addrLen / addresses.length)))
      .join(", ");
  });

  return { connected, connectedText };
}
