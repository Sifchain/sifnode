import { createToken, Token } from "../../entities";

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
  return createToken(1, address, 6, symbol, name);
}

export async function getFakeTokens(): Promise<Map<string, Token>> {
  // gonna load the json and parse the code for all our fake tokens
  const atkJson = require("../../../../chains/ethereum/build/contracts/AliceToken.json");
  const btkJson = require("../../../../chains/ethereum/build/contracts/BobToken.json");

  // Return the tokens parsed as assets
  return new Map([
    ["ATK", parseTruffleJson("AliceToken", "ATK", atkJson)],
    ["BTK", parseTruffleJson("BobToken", "BTK", btkJson)],
  ]);
}
