import { ChainId, createToken, Token } from "../../../entities";

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
  return createToken(symbol, 6, name, ChainId.ETHEREUM, address);
}

// The reason we have to get fake tokens instead of just loading up a
// json list is that everytime truffle compiles we have new addresses
// and we need to keep track of that
export async function getFakeTokens(): Promise<Token[]> {
  // gonna load the json and parse the code for all our fake tokens
  const atkJson = require("../../../../../chains/ethereum/build/contracts/AliceToken.json");
  const btkJson = require("../../../../../chains/ethereum/build/contracts/BobToken.json");

  // Return the tokens parsed as assets
  return [
    parseTruffleJson("AliceToken", "ATK", atkJson),
    parseTruffleJson("BobToken", "BTK", btkJson),
  ];
}
