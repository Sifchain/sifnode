import B from "./B";

test("B utility", () => {
  expect(B(100.12, 2).toString()).toBe("10012");
  expect(B(-100.12, 2).toString()).toBe("-10012");
  expect(B(100.123456789, 10).toString()).toBe("1001234567890");
  expect(B(-100.123456789, 10).toString()).toBe("-1001234567890");
  expect(B("-100.123456789012345678", 18).toString()).toBe(
    "-100123456789012345678",
  );
  expect(B("100.123456789012345678", 18).toString()).toBe(
    "100123456789012345678",
  );
  expect(B("0.00000000000000001", 18).toString()).toBe("10");
  expect(B("1", 18).toString()).toBe("1000000000000000000");
  expect(B("10", 6).toString()).toBe("10000000");
  expect(B("10.000000", 6).toString()).toBe("10000000");
  expect(() => B("-100.12345678.9012345678", 18).toString()).toThrow();
  expect(B("10.000000", 0).toString()).toBe("10");
});
