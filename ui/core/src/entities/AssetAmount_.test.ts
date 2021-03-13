import JSBI from "jsbi";
import { AssetAmount } from "./AssetAmount_";
import { Asset } from "./Asset_";
import { Amount } from "./Amount";
import { Network } from "./Network";

describe("AssetAmount", () => {
  beforeEach(() => {
    // ensure asset is available
    // TODO: Add this to the test utils
    Asset({
      address: "1234568",
      decimals: 18,
      label: "ETH",
      name: "Ethereum",
      network: Network.ETHEREUM,
      symbol: "eth",
      imageUrl: "http://fooo",
    });
  });

  test("Asset values", () => {
    const bal = AssetAmount("eth", "12345678");
    expect(bal.symbol).toBe("eth");
    expect(bal.label).toBe("ETH");
    expect(bal.network).toBe("ethereum");
    expect(bal.imageUrl).toBe("http://fooo");
    expect(bal.decimals).toBe(18);
    expect(bal.address).toBe("1234568");
  });

  test("construction", () => {
    expect(() => AssetAmount("eth", "1234.1234")).toThrow();
  });

  test("Parse to Amount", () => {
    expect(Amount(AssetAmount("eth", "1234")).equalTo(Amount("1234"))).toBe(
      true
    );
  });

  test("#toBigInt", () => {
    // Bigint
    expect(
      JSBI.equal(JSBI.BigInt("1"), AssetAmount("eth", "1").toBigInt())
    ).toBe(true);
  });

  test("#toString", () => {
    expect(AssetAmount("eth", "12345678").toString()).toBe("12345678 ETH");
  });

  test("#add", () => {
    expect(
      AssetAmount("eth", "1000")
        .add(AssetAmount("eth", "1000"))
        .equalTo(AssetAmount("eth", "2000"))
    ).toBe(true);
  });

  describe("#divide", () => {
    test("basic division", () => {
      expect(
        AssetAmount("eth", "10")
          .divide(AssetAmount("eth", "5"))
          .equalTo(AssetAmount("eth", "2"))
      ).toBe(true);
      expect(
        AssetAmount("eth", "30")
          .divide(AssetAmount("eth", "15"))
          .equalTo(AssetAmount("eth", "2"))
      ).toBe(true);
    });

    test("floors remainder", () => {
      expect(
        AssetAmount("eth", "30")
          .divide(AssetAmount("eth", "20"))
          .equalTo(AssetAmount("eth", "1"))
      ).toBe(true);

      expect(
        AssetAmount("eth", "30")
          .divide(AssetAmount("eth", "40"))
          .equalTo(AssetAmount("eth", "0"))
      ).toBe(true);
    });
  });

  test("#equalTo", () => {
    expect(AssetAmount("eth", "1").equalTo(AssetAmount("eth", "1"))).toBe(true);
    expect(AssetAmount("eth", "1").equalTo(AssetAmount("eth", "0"))).toBe(
      false
    );
  });

  test("#greaterThan", () => {
    expect(
      AssetAmount("eth", "100").greaterThan(AssetAmount("eth", "99"))
    ).toBe(true);
    expect(
      AssetAmount("eth", "100").greaterThan(AssetAmount("eth", "100"))
    ).toBe(false);
    expect(
      AssetAmount("eth", "100").greaterThan(AssetAmount("eth", "101"))
    ).toBe(false);
  });

  test("#greaterThanOrEqual", () => {
    expect(
      AssetAmount("eth", "100").greaterThanOrEqual(AssetAmount("eth", "99"))
    ).toBe(true);
    expect(
      AssetAmount("eth", "100").greaterThanOrEqual(AssetAmount("eth", "100"))
    ).toBe(true);
    expect(
      AssetAmount("eth", "100").greaterThanOrEqual(AssetAmount("eth", "101"))
    ).toBe(false);
  });

  test("#lessThan", () => {
    expect(AssetAmount("eth", "100").lessThan(AssetAmount("eth", "99"))).toBe(
      false
    );
    expect(AssetAmount("eth", "100").lessThan(AssetAmount("eth", "100"))).toBe(
      false
    );
    expect(AssetAmount("eth", "100").lessThan(AssetAmount("eth", "101"))).toBe(
      true
    );
  });

  test("#lessThanOrEqual", () => {
    expect(
      AssetAmount("eth", "100").lessThanOrEqual(AssetAmount("eth", "99"))
    ).toBe(false);
    expect(
      AssetAmount("eth", "100").lessThanOrEqual(AssetAmount("eth", "100"))
    ).toBe(true);
    expect(
      AssetAmount("eth", "100").lessThanOrEqual(AssetAmount("eth", "101"))
    ).toBe(true);
  });

  test("#multiply", () => {
    expect(
      AssetAmount("eth", "12345678")
        .multiply(AssetAmount("eth", "10"))
        .equalTo(AssetAmount("eth", "123456780"))
    ).toBe(true);
  });

  test("#sqrt", () => {
    expect(
      AssetAmount("eth", "15241383936")
        .sqrt()
        .equalTo(AssetAmount("eth", "123456"))
    ).toBe(true);

    expect(
      AssetAmount("eth", "15241578750190521")
        .sqrt()
        .equalTo(AssetAmount("eth", "123456789"))
    ).toBe(true);

    // Floor
    expect(
      AssetAmount("eth", "20")
        .sqrt()
        .equalTo(AssetAmount("eth", "4"))
    ).toBe(true);
  });

  test("#subtract", () => {
    expect(
      AssetAmount("eth", "12345678")
        .subtract(AssetAmount("eth", "2345678"))
        .equalTo(AssetAmount("eth", "10000000"))
    ).toBe(true);
  });
});
