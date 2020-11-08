import { Asset, Network, Coin, Token } from "../../../entities";

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
let _memoized: Asset[] | null = null;
export async function getFakeTokens(): Promise<Asset[]> {
  if (_memoized) return _memoized;

  const ETH = Coin({
    symbol: "eth",
    decimals: 18,
    name: "Ethereum",
    network: Network.ETHEREUM,
  });

  // gonna load the json and parse the code for all our fake tokens
  const atkJson = require("../../../../../chains/ethereum/build/contracts/AliceToken.json");
  const btkJson = require("../../../../../chains/ethereum/build/contracts/BobToken.json");
  const ATK = parseTruffleJson("AliceToken", "atk", atkJson);
  const BTK = parseTruffleJson("BobToken", "btk", btkJson);

  // Return the tokens parsed as assets
  _memoized = [ATK, BTK, ETH];
  return _memoized;
}
