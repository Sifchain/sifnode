import { getChains } from "../utils/getChains.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";
import { startChain } from "./startChain.mjs";

export async function startAllChains({
  network,
  home = `/tmp/localnet`,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
}) {
  const chains = getChains({
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
    home,
  });

  return Promise.all(
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
  );

  // if (initRelayer) {
  //   await initRelayers({ chainsProps });
  //   const procs = await startRelayers({ chainsProps });

  //   await Promise.all(
  //     procs.map(async ({ proc }) => {
  //       for await (let chunk of proc.stderr) {
  //         if (chunk.includes("waking up and checking for packets!")) break;
  //       }
  //       proc.kill("SIGINT");
  //     })
  //   );

  //   return;
  // }

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
