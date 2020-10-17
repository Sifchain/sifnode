import { Token } from "../../entities";

export type TokenServiceContext = {
  getSupportedTokens: () => Promise<Token[]>;
};

export default function createWalletService({
  getSupportedTokens,
}: TokenServiceContext) {
  return {
    async getTop20Tokens() {
      return await getSupportedTokens(); // Will be an actually list of the specific tokens we need in order
    },
  };
}
