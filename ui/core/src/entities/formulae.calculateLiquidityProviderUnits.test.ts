import { calculatePoolUnits } from "./formulae";

import tables from "../../../../test/test-tables/pool_units_after_upgrade.json";
import { Amount, IAmount } from "./Amount";

// Use this list to only run specific tests
const filterList: number[] = [];

// The following copies how the backend tests against rounding
// See https://github.com/Sifchain/sifnode/blob/b4a18903319ba3dd5349deb4d6182140e720b163/x/clp/keeper/table_test.go#L14
const BUFFER_PERCENTAGE = "1.21"; // Percentage difference allowable to accommodate rounding done by Big libraries in Go,Python and Javascript

function isAllowable(a: IAmount, b: IAmount) {
  let diffPercentage: IAmount;
  if (a.greaterThanOrEqual(b)) {
    diffPercentage = a.subtract(b).divide(a).multiply("100");
  } else {
    diffPercentage = b.subtract(a).divide(b).multiply("100");
  }

  return !diffPercentage.greaterThanOrEqual(Amount(BUFFER_PERCENTAGE));
}

tables.PoolUnitsAfterUpgrade.forEach(({ r, a, R, A, P, expected }, index) => {
  if (filterList.length === 0 || filterList.includes(index)) {
    test(`#${index} => (r:${r}, a:${a}, R:${R}, A:${A}, P:${P}) => ${expected}`, () => {
      const output = calculatePoolUnits(
        Amount(r),
        Amount(a),
        Amount(R),
        Amount(A),
        Amount(P),
      );

      expect(isAllowable(output, Amount(expected))).toBe(true);
    });
  }
});
