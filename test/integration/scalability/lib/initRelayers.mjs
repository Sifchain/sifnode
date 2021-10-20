import { createRelayer } from "../utils/createRelayer.mjs";
import { createRelayerRegistry } from "../utils/createRelayerRegistry.mjs";
import { send } from "./send.mjs";

export async function initRelayers({
  chainsProps,
  registryFrom = `/tmp/localnet/registry`,
}) {
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
      const relayerHome = `${home}/relayer`;

      await nothrow($`ibc-setup ics20 --home ${relayerHome}`);
    })
  );
}
