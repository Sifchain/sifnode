import { createRequire } from "module";
const require = createRequire(import.meta.url);
const defaultChains = require("../config/chains.json");

export function getChains({
  chains = defaultChains,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
  configPath = `/tmp/localnet/config`,
}) {
  const newChains = { ...chains };

  Object.keys(newChains).forEach((chain, index) => {
    newChains[chain] = {
      ...newChains[chain],
      rpcPort: rpcInitialPort + index,
      p2pPort: p2pInitialPort + index,
      pprofPort: pprofInitialPort + index,
      home: `${configPath}/${chain}/${newChains[chain].chainId}`,
    };
  });

  return newChains;
}
