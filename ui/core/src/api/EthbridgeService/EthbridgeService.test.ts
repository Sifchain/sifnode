import createEthbridgeService from ".";
import localethereumassets from "../../assets.ethereum.localnet.json";
import localsifassets from "../../assets.sifchain.localnet.json";
import { Asset, AssetAmount } from "../../entities";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { getTokenFromSupported } from "../utils/getTokenFromSupported";
import { advanceBlock } from "../utils/advanceBlock";
import createSifService from "../SifService";
import { createWaitForBalance } from "../../test/utils/waitForBalance";

const mnemonic =
  "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow";

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

  test("lock tokens", async () => {
    // get sif balance
    const sifService = createSifService({
      sifApiUrl: "http://localhost:1317",
      sifAddrPrefix: "sif",
      assets: [],
    });

    const sifAccount = "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5";

    await sifService.setPhrase(mnemonic);
    const waitForBalance = createWaitForBalance(sifService);
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
