import { createToken } from "../../entities";

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

export async function getFakeTokens() {
  // gonna load the json and parse the code for all our fake tokens
  const atkJson = require("../../../fixtures/ethereum/build/contracts/AliceToken.json");
  const btkJson = require("../../../fixtures/ethereum/build/contracts/BobToken.json");

  // Return the tokens parsed as assets
  return [
    parseTruffleJson("AliceToken", "ATK", atkJson),
    parseTruffleJson("BobToken", "BTK", btkJson),
  ];
}
