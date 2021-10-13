import { getChainProps } from "../utils/getChainProps.mjs";
import { initRelayers } from "./initRelayers.mjs";
import { startChain } from "./startChain.mjs";
import { startRelayers } from "./startRelayers.mjs";

const chains = require("../config/chains.json");

export async function startAllChains({
  network,
  home = `/tmp/localnet`,
  rpcPort = 11000,
  p2pPort = 12000,
  pprofPort = 13000,
  initRelayer = false,
}) {
  const chainsProps = (
    await Promise.all(
      Object.entries(chains)
        .filter(([_, { disabled = false }]) => disabled === false)
        .map(async ([chain], index) => {
          const chainProps = getChainProps({ chain, network });

          return startChain({
            ...chainProps,
            rpcPort: rpcPort + index,
            p2pPort: p2pPort + index,
            pprofPort: pprofPort + index,
            home: `${home}/${chainProps.chain}/${chainProps.chainId}`,
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

  const procs = await startRelayers({ chainsProps });

  await Promise.all(
    procs.map(async ({ proc }) => {
      for await (let chunk of proc.stderr) {
        console.log(`######`);
        console.log(chunk);
        if (chunk.includes("waking up")) break;
      }
      proc.kill("SIGINT");
    })
  );

  // await Object.values(chainsProps).map(async ({ proc }) => {
  //   proc.kill("SIGINT");
  // });
}
