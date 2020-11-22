import JSBI from "jsbi";

// Convenience method for converting a floating point number with decimals
// to a bigint representation according to a number of decimals
export default function B(num: string | number, dec: number = 18) {
  const numstr = typeof num !== "string" ? num.toFixed(dec) : num;
  const [s, m = "0", huh] = numstr.split(".");
  if (typeof huh !== "undefined") throw new Error("Invalid number string");
  const mm = m.length > dec ? m.slice(0, dec) : m.padEnd(dec, "0");
  const n = [s, mm].join("").replace(/^0+/, "");
  return JSBI.BigInt(n);
}
