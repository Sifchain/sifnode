import { Amount } from "../entities";
import { getMantissaFromDynamicMantissa } from "./format";

const mantissaRange = {
  10000: 2,
  1: 6,
  1000: 4,
  infinity: 0,
};

test("", () => {
  expect(getMantissaFromDynamicMantissa(Amount("500000"), mantissaRange)).toBe(
    0,
  );

  expect(getMantissaFromDynamicMantissa(Amount("10000"), mantissaRange)).toBe(
    0,
  );
  expect(getMantissaFromDynamicMantissa(Amount("9999"), mantissaRange)).toBe(2);
  expect(getMantissaFromDynamicMantissa(Amount("1000"), mantissaRange)).toBe(2);
  expect(getMantissaFromDynamicMantissa(Amount("999"), mantissaRange)).toBe(4);
  expect(getMantissaFromDynamicMantissa(Amount("0.5"), mantissaRange)).toBe(6);
});
