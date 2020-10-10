import { AssetAmount } from "../entities";

function createWalletService() {
  return {
    async getAssetBalances(): Promise<AssetAmount[]> {
      return [];
    },
  };
}

export const walletService = createWalletService();
