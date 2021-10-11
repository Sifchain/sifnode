import { generateRelayerRegistry } from "../utils/generateRelayerRegistry.mjs";
import { send } from "./send.mjs";

export async function initRelayer({
  chainsProps,
  registryFrom = `/tmp/localnet/registry`,
}) {
  await $`mkdir -p ${registryFrom}`;
  const registry = generateRelayerRegistry(chainsProps);
  await fs.writeFile(`${registryFrom}/registry.yaml`, registry);

  const { sifchain: sifChainProps, ...otherChainsProps } = chainsProps;

  const sifSendRequests = [];

  await Promise.all(
    Object.values(otherChainsProps).map(async (otherChainProps) => {
      const { chain, home } = otherChainProps;
      const relayerHome = `${home}/relayer`;
      await nothrow($`mkdir -p ${relayerHome}`);
      await nothrow(
        $`ibc-setup init --home ${relayerHome} --registry-from ${registryFrom} --src ${sifChainProps.chain} --dest ${chain}`
      );
      let addresses = await $`ibc-setup keys list --home ${relayerHome}`;
      addresses = addresses.toString().split("\n");
      const sifChainAddress = addresses
        .find((item) => item.includes(`${sifChainProps.chain}`))
        .replace(`${sifChainProps.chain}: `, ``);
      const otherChainAddress = addresses
        .find((item) => item.includes(`${chain}`))
        .replace(`${chain}: `, ``);
      console.log(`sifChainAddress: ${sifChainAddress}`);
      console.log(`otherChainAddress: ${otherChainAddress}`);

      console.log(sifChainProps);
      sifSendRequests.push({
        ...sifChainProps,
        src: `${sifChainProps.chain}-source`,
        dst: sifChainAddress,
        amount: 10e10,
        node: `tcp://127.0.0.1:${sifChainProps.rpcPort}`,
      });
      console.log(otherChainProps);
      await send({
        ...otherChainProps,
        src: `${chain}-source`,
        dst: otherChainAddress,
        amount: 10e10,
        node: `tcp://127.0.0.1:${otherChainProps.rpcPort}`,
      });
    })
  );

  for await (let sifSendRequest of sifSendRequests) {
    await send(sifSendRequest);
  }

  await sleep(1000);

  await Promise.all(
    Object.values(otherChainsProps).map(async ({ home }) => {
      const relayerHome = `${home}/relayer`;

      await nothrow($`ibc-setup ics20 --home ${relayerHome}`);
      await nothrow($`ibc-relayer start -v --poll 10 --home ${relayerHome}`);
    })
  );
}
