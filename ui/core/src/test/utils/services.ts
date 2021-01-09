// Consolodated place where we can setup testing services

import sifServiceInitializer from "../../api/SifService";
import { TestSifAccount } from "./accounts";
import { getTestingTokens } from "./getTestingToken";

export function createTestSifService(account?: TestSifAccount) {
  const sif = sifServiceInitializer({
    sifApiUrl: "http://localhost:1317",
    sifAddrPrefix: "sif",
    sifWsUrl: "ws://localhost:26657/websocket",
    assets: getTestingTokens(["CATK", "CBTK", "CETH", "ROWAN"]),
  });

  if (account) {
    sif.setPhrase(account.mnemonic);
  }

  return sif;
}
