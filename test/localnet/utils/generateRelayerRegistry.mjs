import { dump } from "js-yaml";

function getChainSettings(chainProps) {
  return {
    [chainProps.chain]: {
      chain_id: chainProps.chainId,
      prefix: chainProps.prefix,
      gas_price: `0.1${chainProps.denom}`,
      hd_path: `m/44'/108'/0'/6'`,
      ics20_port: "transfer",
      rpc: [`http://localhost:${chainProps.rpcPort}`],
    },
  };
}

export function generateRelayerRegistry(chainsProps) {
  return dump({
    version: 1,
    chains: Object.values(chainsProps).reduce(
      (acc, cur) => ({ ...acc, ...getChainSettings(cur) }),
      {}
    ),
  });
}
