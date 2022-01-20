import { $, nothrow } from "zx";
import { send } from "../lib/send.mjs";

export async function createRelayer({
  sifChainProps,
  otherChainProps,
  registryFrom = `/tmp/localnet/config/registry`,
}) {
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

  await send({
    ...otherChainProps,
    src: `${chain}-source`,
    dst: otherChainAddress,
    amount: 10e10,
    node: `tcp://127.0.0.1:${otherChainProps.rpcPort}`,
  });

  return {
    sifChainAddress,
    otherChainAddress,
    sifSendRequest: {
      ...sifChainProps,
      src: `${sifChainProps.chain}-source`,
      dst: sifChainAddress,
      amount: 10e10,
      node: `tcp://127.0.0.1:${sifChainProps.rpcPort}`,
    },
  };
}
