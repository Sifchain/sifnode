export function isSupportedEVMChain(chainId?: string) {
  if (!chainId) return false;
  // List of supported EVM chainIds
  const supportedEVMChainIds = [
    "0x1", // 1 Mainnet
    "0x3", // 3 Ropsten
    "0x539", // 1337 Ganache/Hardhat
  ];

  return supportedEVMChainIds.includes(chainId);
}
