import createEthbridgeService from ".";
import localethereumassets from "../../assets.ethereum.localnet.json";
import localsifassets from "../../assets.sifchain.localnet.json";
import { Asset, AssetAmount } from "../../entities";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { getTokenFromSupported } from "../utils/getTokenFromSupported";
import { advanceBlock } from "../utils/advanceBlock";
import createSifService, { ISifService } from "../SifService";

import Web3 from "web3";
import JSBI from "jsbi";

const mnemonic =
  "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow";
const sleep = (ms: number) => new Promise((done) => setTimeout(done, ms));
const balanceWaiter = (sifService: ISifService) => {
  async function getBalance(symbol: string) {
    return (
      await sifService.getBalance("sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd")
    ).find((bal) => bal.asset.symbol === symbol);
  }

  return async function checkBalance(
    symbol: string,
    expectedAmount: string,
    account: string,
    maxTries = 100
  ) {
    for (let i = 0; i < maxTries; i++) {
      await sleep(1000);

      const newBalance = (await sifService.getBalance(account)).find(
        (bal) => bal.asset.symbol === symbol
      );

      if (newBalance?.amount.toString() === expectedAmount) {
        return newBalance;
      }
    }
    throw new Error(`Balance of ${expectedAmount} was never realised`);
  };
};

describe("PeggyService", () => {
  let EthbridgeService: ReturnType<typeof createEthbridgeService>;

  let ETH: Asset;

  beforeEach(async () => {
    require("@openzeppelin/test-helpers/configure")({
      provider: await getWeb3Provider(),
    });

    ETH = getTokenFromSupported(
      [...localethereumassets.assets, ...localsifassets.assets],
      "ETH"
    );

    EthbridgeService = createEthbridgeService({
      bridgebankContractAddress: "0xf204a4Ef082f5c04bB89F7D5E6568B796096735a",
      getWeb3Provider,
    });
  });

  test.only("lock tokens", async () => {
    // get sif balance
    const sifService = createSifService({
      sifApiUrl: "http://localhost:1317",
      sifAddrPrefix: "sif",
      assets: [],
    });

    const sifAccount = "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5";

    await sifService.setPhrase(mnemonic);
    const waitForBalance = balanceWaiter(sifService);
    await waitForBalance("ceth", "1000000000", sifAccount);

    await new Promise<void>(async (done) => {
      EthbridgeService.lock(sifAccount, AssetAmount(ETH, "2"), 10)
        .onTxEvent((evt) => {
          console.log(evt);
        })
        .onComplete(async () => {
          // Not testing balances because we have no
          // way to correlate against transaction
          const balance = await waitForBalance(
            "ceth",
            "2000000001000000000",
            sifAccount
          );

          console.log("balance.amount:", balance.amount);
          done();
        })
        .onError((err) => {
          throw err.payload;
        });
      advanceBlock(200);
    });
  });
});
