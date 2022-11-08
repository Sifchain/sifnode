export function pickChains({ chain }) {
  if (!chain) {
    throw new Error("chain not defined");
  }

  const chains = chain.split(",");
  if (chains.length === 0) {
    throw new Error("chains is empty");
  }
  return chains;
}
