import { calculatePoolUnits } from "./formulae";
import { Fraction } from "./fraction/Fraction";

import tables from "../../../../test/test-tables/pool_units.json";

// Use this list to only run specific tests
const filterList: number[] = [];

tables.PoolUnits.forEach(({ r, a, R, A, P, expected }, index) => {
  if (filterList.length === 0 || filterList.includes(index)) {
    test(`#${index} => (r:${r}, a:${a}, R:${R}, A:${A}, P:${P}) => ${expected}`, () => {
      const output = calculatePoolUnits(
        new Fraction(r),
        new Fraction(a),
        new Fraction(R),
        new Fraction(A),
        new Fraction(P)
      );

      expect(output.toFixed(0)).toBe(expected);
    });
  }
});
