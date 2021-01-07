import createEthbridgeService from ".";
import { Asset, AssetAmount } from "../../entities";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";

import { advanceBlock } from "../../test/utils/advanceBlock";
import { createWaitForBalance } from "../../test/utils/waitForBalance";
import { akasha } from "../../test/utils/accounts";
import { createTestSifService } from "../../test/utils/services";
import { getTestingToken } from "../../test/utils/getTestingToken";

describe("PeggyService", () => {
  let EthbridgeService: ReturnType<typeof createEthbridgeService>;

  let ETH: Asset;

  beforeEach(async () => {
    ETH = getTestingToken("ETH");

    EthbridgeService = createEthbridgeService({
      bridgebankContractAddress: "0xf204a4Ef082f5c04bB89F7D5E6568B796096735a",
      getWeb3Provider,
    });
  });

  test("lock tokens", async () => {
    // get sif balance
    const sifService = createTestSifService(akasha);

    const waitForBalance = createWaitForBalance(sifService);
    await waitForBalance("ceth", "1000000000", akasha.address);

    await new Promise<void>(async (done) => {
      EthbridgeService.lock(akasha.address, AssetAmount(ETH, "2"), 10)
        .onTxEvent((evt) => {
          console.log(evt);
        })
        .onComplete(async () => {
          // Not testing balances because we have no
          // way to correlate against transaction
          await waitForBalance("ceth", "2000000001000000000", akasha.address);
          done();
        })
        .onError((err) => {
          throw err.payload;
        });
      advanceBlock(200);
    });
  });
});
