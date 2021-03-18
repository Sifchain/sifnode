import JSBI from "jsbi";
import { Amount } from "./Amount";

describe("Amount", () => {
  test("construction", () => {
    expect(() => Amount("1234.1234")).toThrow();
  });

  test("#toBigInt", () => {
    // Bigint
    expect(JSBI.equal(JSBI.BigInt("1"), Amount("1").toBigInt())).toBe(true);

    // Supports Negative Numbers
    expect(JSBI.equal(JSBI.BigInt("-1234"), Amount("-1234").toBigInt())).toBe(
      true,
    );

    expect(JSBI.equal(JSBI.BigInt("1"), Amount("2").toBigInt())).toBe(false);
  });

  test("#toString", () => {
    expect(Amount("12345678").toString()).toBe("12345678");
  });

  test("#add", () => {
    expect(Amount("1000").add(Amount("1000")).equalTo(Amount("2000"))).toBe(
      true,
    );
  });

  describe("#divide", () => {
    test("basic division", () => {
      expect(Amount("10").divide(Amount("5")).equalTo(Amount("2"))).toBe(true);
      expect(Amount("30").divide(Amount("15")).equalTo(Amount("2"))).toBe(true);
    });

    test("floors remainder", () => {
      expect(Amount("30").divide(Amount("20")).equalTo(Amount("1"))).toBe(true);

      expect(Amount("30").divide(Amount("40")).equalTo(Amount("0"))).toBe(true);
    });
  });

  test("#equalTo", () => {
    expect(Amount("1").equalTo(Amount("1"))).toBe(true);
    expect(Amount("1").equalTo(Amount("0"))).toBe(false);
  });

  test("#greaterThan", () => {
    expect(Amount("100").greaterThan(Amount("99"))).toBe(true);
    expect(Amount("100").greaterThan(Amount("100"))).toBe(false);
    expect(Amount("100").greaterThan(Amount("101"))).toBe(false);
  });

  test("#greaterThanOrEqual", () => {
    expect(Amount("100").greaterThanOrEqual(Amount("99"))).toBe(true);
    expect(Amount("100").greaterThanOrEqual(Amount("100"))).toBe(true);
    expect(Amount("100").greaterThanOrEqual(Amount("101"))).toBe(false);
  });

  test("#lessThan", () => {
    expect(Amount("100").lessThan(Amount("99"))).toBe(false);
    expect(Amount("100").lessThan(Amount("100"))).toBe(false);
    expect(Amount("100").lessThan(Amount("101"))).toBe(true);
  });

  test("#lessThanOrEqual", () => {
    expect(Amount("100").lessThanOrEqual(Amount("99"))).toBe(false);
    expect(Amount("100").lessThanOrEqual(Amount("100"))).toBe(true);
    expect(Amount("100").lessThanOrEqual(Amount("101"))).toBe(true);
  });

  test("#multiply", () => {
    expect(
      Amount("12345678").multiply(Amount("10")).equalTo(Amount("123456780")),
    ).toBe(true);
  });

  test("#sqrt", () => {
    expect(Amount("15241383936").sqrt().equalTo(Amount("123456"))).toBe(true);

    expect(
      Amount("15241578750190521").sqrt().equalTo(Amount("123456789")),
    ).toBe(true);

    // Floor
    expect(Amount("20").sqrt().equalTo(Amount("4"))).toBe(true);
  });

  test("#subtract", () => {
    expect(
      Amount("12345678")
        .subtract(Amount("2345678"))
        .equalTo(Amount("10000000")),
    ).toBe(true);
  });
});
