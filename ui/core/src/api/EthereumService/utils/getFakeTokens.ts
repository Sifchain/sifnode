import { Asset, Network, Coin, Token } from "../../../entities";
import { getSupportedTokens } from "./getSupportedTokens";

// Parse Truffle json for the most recent address
function parseTruffleJson(
  name: string,
  symbol: string,
  jobj: {
    networks: { [s: string]: { address: string } };
  }
) {
  const entries = Object.entries(jobj.networks);
  const [last] = entries.slice(-1);
  const { address } = last[1];
  return Token({
    symbol,
    decimals: 6,
    name,
    network: Network.ETHEREUM,
    address,
  });
}

// The reason we have to get fake tokens instead of just loading up a
// json list is that everytime truffle compiles we have new addresses
// and we need to keep track of that
export async function getFakeTokens(): Promise<Token[]> {
  // add real tokens for testing
  const realTokens = await getSupportedTokens();

  // gonna load the json and parse the code for all our fake tokens
  const atkJson = require("../../../../../chains/ethereum/build/contracts/AliceToken.json");
  const btkJson = require("../../../../../chains/ethereum/build/contracts/BobToken.json");

  // Return the tokens parsed as assets
  return [
    parseTruffleJson("AliceToken", "atk", atkJson),
    parseTruffleJson("BobToken", "btk", btkJson),
    ...realTokens,
  ];
}

export async function getFakeAssets(): Promise<Asset[]> {
  const ETH = Coin({
    symbol: "ETH",
    decimals: 18,
    name: "Ethereum",
    network: Network.ETHEREUM,
  });
  const RWN = Coin({
    symbol: "nametoken",
    decimals: 6,
    name: "nametoken",
    network: Network.SIFCHAIN,
  });

  return [ETH, RWN];
}
