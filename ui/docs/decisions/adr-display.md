// import DisplayAmount from "@core/DisplayAmount";

interface Asset {} // Re-used
interface Amount {} // Re-used - Rudi implmentation

enum Type {
  fixed,
  absolute,
  significant,
}

interface IOptions {
  precision?: number; // Sets the number of decimals to show, defaults to asset precision unless this is set
  floor?: boolean; // Should we round up or down (depending on significant digit)- WIP
  commas?: boolean; // 1000.1000 or 1,000.1000 -> // https://prowritingaid.com/grammar/1008091/When-should-I-use-a-comma-to-separate-numbers
  padding?: number; // If we want to always render things so they are aligned, sometimes you have to add 0's
  // e.g. if we set padding=3 1000.1000
  //                          1000.0001

  meaningfulDollar?: number; // I need to give this some thought to explain, but sometimes you just want to show enough decimals that makes a dollar amount make sense. e.g. If you have 1.123112314 btc ~= 55,517.87, most of the information you need is just contained in perhaps just 1.123 but if you only have 0.0001213 btc ~= 7.72 then you would have to show at least 0.00012. This won't have to be built initially but might show up eventually in product designs. The price is available inside of Asset, so this is probably easy to do. As above, needs more thinking though.
}

// Unfortunately I can't make all options available to all types
// I need some ideas on how to map this out e.g. type: absolute shouldn't really allow for fixed digits
// but type: fixed means that you should be able to change the number of digits

function DisplayAmount(
  asset: Asset,
  amount: Amount,
  type: Type,
  options?: IOptions
): string {
  return "1,000";
}

const asset = {}; // Imagine this is a real asset
const amount = {}; // Imagine these are JSBIS

let balance = DisplayAmount(asset, amount, Type.absolute);
// e.g. amount = 0.000000000000000001 returns 0.000000000000000001

balance = DisplayAmount(asset, amount, Type.absolute, { commas: true });
// e.g. amount = 1000.000000000000000001 returns 1,000.000000000000000001

balance = DisplayAmount(asset, amount, Type.fixed, { precision: 3, padding: 2 });
// e.g. amount = 1000.001 returns 1000.00100;
