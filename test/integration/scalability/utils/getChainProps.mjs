import { createRequire } from "module";
const require = createRequire(import.meta.url);
const defaultChains = require("../config/chains.json");

export function getChainProps({
  chains = defaultChains,
  chain,
  network,
  type,
  ...rest
}) {
  if (chain && !chains[chain])
    throw new Error("this chain name is not defined within chains.json file");

  Object.keys(rest).forEach(
    (key) => rest[key] === undefined && delete rest[key]
  );

  let result = {};

  if (chain) result.chain = chain;
  if (network) result.network = network;
  if (type) result.type = type;

  result = {
    ...result,
    ...(chains[chain] || {}),
    ...((chains[chain] || {})[network] || {}),
    ...(((chains[chain] || {})[network] || {})[type] || {}),
    ...rest,
  };

  return result;
}
