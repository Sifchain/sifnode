import { useCore } from "@/hooks/useCore";
import { computed } from "@vue/reactivity";

const shorten = (num: number) => (str: string) => str.slice(0, num) + "...";

export function useWalletButton(props?: {
  connectType?: "connectToAny" | "connectToAll" | "connectToSif";
  addrLen?: number;
}) {
  const connectType = props?.connectType || "connectToAny";
  const { store } = useCore();

  const connectedToEth = computed(() => store.wallet.eth.isConnected);

  const connectedToSif = computed(() => store.wallet.sif.isConnected);

  const connected = computed(() => {
    if (connectType === "connectToAny") {
      return connectedToSif.value || connectedToEth.value;
    }

    if (connectType === "connectToAll") {
      return connectedToSif.value && connectedToEth.value;
    }

    if (connectType === "connectToSif") {
      return connectedToSif.value;
    }
  });

  const connectedText = computed(() => {
    const addresses = [
      store.wallet.eth.address,
      store.wallet.sif.address,
    ].filter(Boolean);

    return addresses.length > 0 ? "Connected" : "Connected";
  });

  const connectCta = computed(() => {
    if (!(store.wallet.eth.isConnected || store.wallet.sif.isConnected)) {
      return "Connect Wallet";
    }
    if (!store.wallet.sif.isConnected) {
      return "Connect Sifchain Wallet";
    }
    if (!store.wallet.eth.isConnected) {
      return "Connect Ethereum Wallet";
    }
  });

  return {
    connected,
    connectedToEth,
    connectedToSif,
    connectedText,
    connectCta,
  };
}
