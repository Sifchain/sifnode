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
  const atkJson = require("./AliceToken");
  const btkJson = require("./BobToken");
  const ATK = parseTruffleJson("AliceToken", "atk", atkJson);
  const BTK = parseTruffleJson("BobToken", "btk", btkJson);

  const CBTK = Coin({
    symbol: "cbtk",
    decimals: 18,
    name: "Banana",
    network: Network.SIFCHAIN,
  });

  const CATK = Coin({
    symbol: "catk",
    decimals: 18,
    name: "Apple",
    network: Network.SIFCHAIN,
  });

  const CETH = Coin({
    symbol: "ceth",
    decimals: 18,
    name: "Ethereum",
    network: Network.SIFCHAIN,
  });

  const RWN = Coin({
    symbol: "crwn",
    decimals: 18,
    name: "Rowan",
    network: Network.SIFCHAIN,
  });

  // Return the tokens parsed as assets
  _memoized = [ATK, BTK, ETH, CBTK, CATK, CETH, RWN];
  return _memoized;
}
