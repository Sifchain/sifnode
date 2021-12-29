import { getChainProps } from "../utils/getChainProps.mjs";

export function getChainsProps({ chains, network }) {
  return Object.entries(chains)
    .filter(([_, { disabled = false }]) => disabled === false)
    .map(([chain, chainProps]) =>
      getChainProps({ chain, network, ...chainProps })
    )
    .reduce((acc, cur) => ({ ...acc, [cur.chain]: cur }), {});
}
