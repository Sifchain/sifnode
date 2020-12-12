import { ATK } from "../../constants";
import { AssetAmount } from "../../entities";
import create from "./PeggyService";
describe("PeggyService", () => {
  // We are going to test this as a mock implementation
  // These tests may have to change to be integration tests
  // at a later point

  test("lock", async () => {
    const events: any[] = [];
    const service = create();

    expect(events).toEqual([]);

    await new Promise<void>((resolve) => {
      service
        .lock("sif12345876512341234", AssetAmount(ATK, "10000"))
        .onTxEvent((e) => events.push(e))
        .onComplete(() => resolve());
    });

    expect(events.map((e) => e.type)).toEqual([
      "EthTxInitiated",
      ...Array.from(Array(30).keys()).map(() => "EthConfCountChanged"),
      "EthTxConfirmed",
      "SifTxInitiated",
      "SifTxConfirmed",
      "Complete",
    ]);
  });

  test("error", async () => {
    const events: any[] = [];
    const service = create();

    expect(events).toEqual([]);

    await new Promise<void>((resolve) => {
      service
        .lock("sif12345876512341234", AssetAmount(ATK, "100"))
        .onError((e) => {
          events.push(e);
          resolve();
        });
    });

    expect(events).toEqual([
      {
        txHash: "abcd1234", // temp
        type: "Error",
        payload: "Boom!",
      },
    ]);
  });

  test("burn", async () => {
    const events: any[] = [];
    const service = create();

    expect(events).toEqual([]);

    await new Promise<void>((resolve) => {
      service
        .burn("sif12345876512341234", AssetAmount(ATK, "10000"))
        .onTxEvent((e) => events.push(e))
        .onComplete(() => resolve());
    });

    expect(events.map((e) => e.type)).toEqual([
      "SifTxInitiated",
      ...Array.from(Array(30).keys()).map(() => "SifConfCountChanged"),
      "SifTxConfirmed",
      "EthTxInitiated",
      "EthTxConfirmed",
      "Complete",
    ]);
  });
});
