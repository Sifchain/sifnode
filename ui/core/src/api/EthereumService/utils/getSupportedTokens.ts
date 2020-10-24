import { Network, Token } from "../../../entities";
type TokenData = {
  id: string;
  name: string;
  symbol: string;
  image: {
    thumb: string;
    small: string;
    large: string;
  };
  contract_address: string;
  asset_platform_id: string;
  market_cap_rank: number;
  decimals: number;
};

export async function getSupportedTokens(): Promise<Token[]> {
  const tokens: TokenData[] = require("../../../../data/ethereum_tokens.json");
  return tokens.map((t) =>
    Token({
      address: t.contract_address,
      decimals: t.decimals,
      name: t.name,
      network: Network.ETHEREUM,
      symbol: t.symbol.toLowerCase(),
      imageUrl: t.image.small,
    })
  );
}
