import createEthbridgeService from ".";
import localethereumassets from "../../assets.ethereum.localnet.json";
import localsifassets from "../../assets.sifchain.localnet.json";
import { Asset, AssetAmount } from "../../entities";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";
import { getTokenFromSupported } from "../utils/getTokenFromSupported";
import { advanceBlock } from "../utils/advanceBlock";

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

  test("lock tokens", (done) => {
    EthbridgeService.lock(
      "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5",
      AssetAmount(ETH, "0.0002"),
      5
    )
      .onComplete(async () => {
        // Not testing balances because we have no
        // way to correlate against transaction
        done();
      })
      .onError((err) => {
        throw err.payload;
      });

    advanceBlock(7);
  });
});
