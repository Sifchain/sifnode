// This test must be run alongside one ganache and one sifnode instance

import createPeggyService, { IPeggyService } from ".";
import { getWeb3Provider } from "../../test/utils/getWeb3Provider";

describe("PeggyService", () => {
  let PeggyService: IPeggyService;

  beforeEach(async () => {
    PeggyService = createPeggyService({
      sifApiUrl: "",
      bridgeBankContractAddress: "",
      getWeb3Provider,
    });
  });

  test("pegging tokens from ethereum", async () => {});
});
