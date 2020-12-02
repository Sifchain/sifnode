import { Asset, Network, Coin, Token } from "../../../entities";

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

  const ATK = Token({
    symbol: "atk",
    decimals: 6,
    name: "AliceToken",
    network: Network.ETHEREUM,
    address: "0xbaAA2a3237035A2c7fA2A33c76B44a8C6Fe18e87", // NOTE: address dependent on ganache seed and migration order
  });
  const BTK = Token({
    symbol: "btk",
    decimals: 6,
    name: "BobToken",
    network: Network.ETHEREUM,
    address: "0x13274Fe19C0178208bCbee397af8167A7be27f6f", // NOTE: address dependent on ganache seed and migration order
  });

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
