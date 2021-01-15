// Consolodated place where we can setup testing services

import sifServiceInitializer from "../../api/SifService";
import { KeplrChainConfig } from "../../utils/parseConfig";
import { TestSifAccount } from "./accounts";
import { getTestingTokens } from "./getTestingToken";

export function createTestSifService(account?: TestSifAccount) {
  const sif = sifServiceInitializer({
    sifApiUrl: "http://localhost:1317",
    sifAddrPrefix: "sif",
    sifWsUrl: "ws://localhost:26657/websocket",
    assets: getTestingTokens(["CATK", "CBTK", "CETH", "ROWAN"]),
    keplrChainConfig: {} as KeplrChainConfig,
  });

  if (account) {
    sif.setPhrase(account.mnemonic);
  }

  return sif;
}
