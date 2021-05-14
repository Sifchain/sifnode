// Consolodated place where we can setup testing services

import sifServiceInitializer from "../../services/SifService";
import ethServiceInitializer from "../../services/EthereumService";
import { KeplrChainConfig } from "../../utils/parseConfig";
import { TestSifAccount, TestEthAccount } from "./accounts";
import { getTestingTokens } from "./getTestingToken";
import { getWeb3Provider } from "./getWeb3Provider";

export async function createTestSifService(account?: TestSifAccount) {
  const sif = sifServiceInitializer({
    sifApiUrl: "http://localhost:1317",
    sifAddrPrefix: "sif",
    sifWsUrl: "ws://localhost:26657/websocket",
    sifRpcUrl: "http://localhost:26657",
    assets: getTestingTokens(["CATK", "CBTK", "CETH", "ROWAN"]),
    keplrChainConfig: {} as KeplrChainConfig,
  });

  if (account) {
    console.log("logging in to account with: " + account.mnemonic);
    await sif.setPhrase(account.mnemonic);
  }

  return sif;
}

export async function createTestEthService() {
  const eth = ethServiceInitializer({
    assets: getTestingTokens(["ATK", "BTK", "ETH", "EROWAN"]),
    getWeb3Provider,
  });
  console.log("Connecting to eth service");
  await eth.connect();
  console.log("Finished connecting to eth service");
  return eth;
}
