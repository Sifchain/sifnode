import { createRelayer } from "../utils/createRelayer.mjs";
import { createRelayerRegistry } from "../utils/createRelayerRegistry.mjs";
import { setupRelayerChannelIds } from "../utils/setupRelayerChannelIds.mjs";
import { getChains } from "../utils/getChains.mjs";
import { getChainsProps } from "../utils/getChainsProps.mjs";
import { send } from "./send.mjs";
import { sleep } from "zx";

export async function initAllRelayers({
  network,
  configPath = `/tmp/localnet/config`,
  registryFrom = `/tmp/localnet/config/registry`,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
}) {
  // 0) retrieve chains + metadata
  const chains = getChains({
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
    configPath,
  });
  const chainsProps = getChainsProps({ chains, network });
  const { sifchain: sifChainProps, ...otherChainsProps } = chainsProps;

  // 1) create global registry for relayers
  await createRelayerRegistry({ chainsProps, registryFrom });

  // 2) create relayer for each single chains connecting to sifchain
  const createdRelayers = await Promise.all(
    Object.values(otherChainsProps).map(async (otherChainProps) => {
      return createRelayer({ sifChainProps, otherChainProps, registryFrom });
    })
  );

  // 3) fund all relayer addresses
  for await (let createdRelayer of createdRelayers) {
    await send(createdRelayer.sifSendRequest);
  }

  // 4) wait
  await sleep(1000);

  // 5) generate channel IDs
  await Promise.all(
    Object.values(otherChainsProps).map(async ({ home }) => {
      await setupRelayerChannelIds({ home });
    })
  );
}
