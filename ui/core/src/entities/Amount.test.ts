import JSBI from "jsbi";
import { Amount } from "./Amount";

describe("Amount", () => {
  test("construction from decimal", () => {
    expect(Amount("1.5").toString()).toBe("1.500000000000000000");
  });

  test("construction from negative decimal", () => {
    expect(Amount("-1.5").toString()).toBe("-1.500000000000000000");
  });

  test("construction from negative integer", () => {
    expect(Amount("-5").toString()).toBe("-5.000000000000000000");
  });

  test("construction from decimal with no leading zero", () => {
    expect(Amount(".5").toString()).toBe("0.500000000000000000");
  });

  test("construction from integer", () => {
    expect(Amount("15").toString()).toBe("15.000000000000000000");
  });

  test("construction with garbage input should throw", () => {
    expect(() => {
      Amount("1.5.34");
    }).toThrow();
  });

  test("construction with garbage input should throw", () => {
    expect(() => {
      Amount("..534");
    }).toThrow();
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
    expect(Amount("12345678").toString()).toBe("12345678.000000000000000000");
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

    test("bankers rounding", () => {
      expect(Amount("30").divide(Amount("20")).toBigInt().toString()).toBe("2");
      expect(Amount("30").divide(Amount("40")).toBigInt().toString()).toBe("1");
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
    expect(Amount("15241383936").sqrt().toBigInt().toString()).toBe("123456");

    expect(Amount("15241578750190521").sqrt().toBigInt().toString()).toBe(
      "123456789",
    );

    // Floor
    expect(Amount("20").sqrt().toString()).toBe("4.472135954999579393");
  });

  test("#expInt", () => {
    expect(Amount("2").expInt("18").toBigInt().toString()).toBe(
      "2000000000000000000",
    );
    expect(Amount("2").expInt("-18").toString()).toBe("0.000000000000000002");
  });

  test("#powerInt", () => {
    expect(Amount("2").powerInt("-3").toString()).toBe("0.125000000000000000");
    expect(Amount("2").powerInt("-2").toString()).toBe("0.250000000000000000");
    expect(Amount("2").powerInt("-1").toString()).toBe("0.500000000000000000");
    expect(Amount("2").powerInt("0").toBigInt().toString()).toBe("1");
    expect(Amount("2").powerInt("1").toBigInt().toString()).toBe("2");
    expect(Amount("2").powerInt("6").toBigInt().toString()).toBe("64");
    expect(Amount("2").powerInt("8").toBigInt().toString()).toBe("256");
  });

  test("#subtract", () => {
    expect(
      Amount("12345678")
        .subtract(Amount("2345678"))
        .equalTo(Amount("10000000")),
    ).toBe(true);
  });
});
