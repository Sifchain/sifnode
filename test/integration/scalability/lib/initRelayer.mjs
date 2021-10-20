import { createRelayer } from "../utils/createRelayer.mjs";
import { createRelayerRegistry } from "../utils/createRelayerRegistry.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";
import { send } from "./send.mjs";

export async function initRelayer(props) {
  const {
    chainProps: otherChainProps,
    registryFrom = `/tmp/localnet/registry`,
  } = props;
  const {
    rpcPort = 26657,
    p2pPort = 26656,
    pprofPort = 6060,
    home = `/tmp/localnet/${otherChainProps.chain}/${otherChainProps.chainId}`,
  } = otherChainProps;

  // 0) retrieve sifchain props
  const sifChainProps = getChainProps({ chain: "sifchain" });

  // 1) create global registry for relayers
  await createRelayerRegistry({
    chainsProps: [sifChainProps, otherChainProps],
    registryFrom,
  });

  // 2) create relayer for pair of chain
  const createdRelayer = await createRelayer({
    sifChainProps,
    otherChainProps: { ...otherChainProps, rpcPort, p2pPort, pprofPort, home },
    registryFrom,
  });

  // 3) fund all relayer addresses
  await send(createdRelayer.sifSendRequest);

  // 4) wait
  await sleep(1000);

  // 5) generate channel IDs
  await setupRelayerChannelIds({ home });
}
