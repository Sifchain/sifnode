export function ensureSifAddress(address: string) {
  if (address.length !== 42) throw "Address not valid (length). Fail"; // this is simple check, limited to default address type (check bech32);
  if (!address.match(/^sif/)) throw "Address not valid (format). Fail"; // this is simple check, limited to default address type (check bech32);
  // TODO: add invariant address starts with "sif" (double check this is correct)
  return address;
}
