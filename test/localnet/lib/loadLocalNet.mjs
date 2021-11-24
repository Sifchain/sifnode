import { $ } from "zx";
import pEvent from "p-event";
import { extractArchive } from "../utils/extractArchive.mjs";
import { startAllChains } from "./startAllChains.mjs";
import { startAllRelayers } from "./startAllRelayers.mjs";

export async function loadLocalNet({
  configPath = "/tmp/localnet/config",
  archivePath = "/tmp/localnet/config.tbz",
  network,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
}) {
  // 0) first make sure we start from an empty targeted folder
  await $`rm -rf ${configPath}`;

  // 1) extract the snapshot archive into the targeted folder
  await extractArchive({ archivePath, configPath });

  try {
    // 2) start all IBC chains
    const chainsProps = await startAllChains({
      network,
      configPath,
      rpcInitialPort,
      p2pInitialPort,
      pprofInitialPort,
    });

    await Promise.all(
      chainsProps.map(async ({ proc }) => {
        const asyncIterator = pEvent.iterator(proc.stderr, "data", {
          resolutionEvents: ["finish"],
        });
        for await (let chunk of asyncIterator) {
          if (chunk.includes("Starting RPC HTTP server")) break;
        }
      })
    );

    // 3) start all IBC relayers
    const relayersProps = await startAllRelayers({
      network,
      configPath,
      rpcInitialPort,
      p2pInitialPort,
      pprofInitialPort,
    });

    // await Promise.all(
    //   relayersProps.map(async ({ proc }) => {
    //     const asyncIterator = pEvent.iterator(proc.stdout, "data", {
    //       resolutionEvents: ["finish"],
    //     });
    //     for await (let chunk of asyncIterator) {
    //       if (chunk.includes("Use last queried heights")) break;
    //     }
    //   })
    // );

    return { chainsProps, relayersProps };
  } catch (err) {
    // ignore
  }
}
