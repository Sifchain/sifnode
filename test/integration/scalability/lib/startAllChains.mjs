import { getChains } from "../utils/getChains.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";
import { initRelayers } from "./initRelayers.mjs";
import { startChain } from "./startChain.mjs";
import { startRelayers } from "./startRelayers.mjs";

export async function startAllChains({
  network,
  home = `/tmp/localnet`,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
  initRelayer = false,
}) {
  const chains = getChains({
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
    home,
  });

  const chainsProps = (
    await Promise.all(
      Object.entries(chains)
        .filter(([_, { disabled = false }]) => disabled === false)
        .map(async ([chain, chainProps]) => {
          const newChainProps = getChainProps({
            chain,
            network,
            ...chainProps,
          });
          return startChain({
            ...newChainProps,
          });
        })
    )
  ).reduce((acc, cur) => ({ ...acc, [cur.chain]: cur }), {});

  if (initRelayer) {
    await initRelayers({ chainsProps });
    const procs = await startRelayers({ chainsProps });

    await Promise.all(
      procs.map(async ({ proc }) => {
        for await (let chunk of proc.stderr) {
          if (chunk.includes("waking up and checking for packets!")) break;
        }
        proc.kill("SIGINT");
      })
    );

    return;
  }

  // const procs = await startRelayers({ chainsProps });

  // await Promise.all(
  //   procs.map(async ({ proc }) => {
  //     for await (let chunk of proc.stderr) {
  //       console.log(`######`);
  //       console.log(chunk);
  //       if (chunk.includes("waking up")) break;
  //     }
  //     proc.kill("SIGINT");
  //   })
  // );

  // await Object.values(chainsProps).map(async ({ proc }) => {
  //   proc.kill("SIGINT");
  // });
}
