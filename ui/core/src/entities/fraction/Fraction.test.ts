import { Fraction } from "./Fraction";

test("it should be able to handle whole integars", () => {
  const f = new Fraction("100");
  expect(f.toFixed(2)).toBe("100.00");
});

test("it should be able to handle negative integars", () => {
  const f = new Fraction("-10015", "100");
  expect(f.toFixed(2)).toBe("-100.15");
});
