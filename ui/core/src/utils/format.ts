import { AssetAmount, IAssetAmount } from "../entities/AssetAmount";

export default function format(amount: IAssetAmount): string {
  return amount.toFormatted();
}

// maybe 80% of the time when building out the UX there will be named types
// 'short-amount', 'short-meaning-ful-balance' etc
// format(aAmount, 'short-amount');
// vs
// format(aAmount, { decimals: 2, commas: true })
// this might be over thinking, maybe it isn't that bad having explicit values
// over and over again in calls to format
// and/or keep format explicit and rely on components to handle abstractions

/*
perhaps another export formatF etc

format(aAmount, 'large-price') {
  return format(aAmount, { decimals: 8, commas: true })
}

*/
