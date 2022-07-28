export function getFunds(tokens, separator = ",") {
  const funds = `${tokens
    .map((token) => `999000000000000000000000000000000${token}`)
    .join(separator)}${separator}999000000000000000000000000000000stake`;
  return funds;
}
