import { useCore } from "@/hooks/useCore";
import { computed } from "@vue/reactivity";

export function useWalletButton(props?: {
  connectType?: "connectToAny" | "connectToAll" | "connectToSif";
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
    connectCta,
  };
}
