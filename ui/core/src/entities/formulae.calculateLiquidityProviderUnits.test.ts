import { calculatePoolUnits } from "./formulae";

import tables from "../../../../test/test-tables/pool_units.json";
import { Amount } from "./Amount";

// Use this list to only run specific tests
const filterList: number[] = [];

tables.PoolUnits.forEach(({ r, a, R, A, P, expected }, index) => {
  if (filterList.length === 0 || filterList.includes(index)) {
    test(`#${index} => (r:${r}, a:${a}, R:${R}, A:${A}, P:${P}) => ${expected}`, () => {
      const output = calculatePoolUnits(
        Amount(r),
        Amount(a),
        Amount(R),
        Amount(A),
        Amount(P),
      );

      expect(output.toString()).toBe(expected);
    });
  }
});
