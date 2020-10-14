import { Token } from "../entities";

export type TokenServiceContext = {
  getSupportedTokens: () => Promise<Token[]>;
};

export default function createWalletService({
  getSupportedTokens,
}: TokenServiceContext) {
  return {
    async getTopERC20Tokens({ limit }: { limit: number }): Promise<Token[]> {
      const tokens = await getSupportedTokens();
      // TODO: order the tokens by market cap
      return tokens;
    },
  };
}
