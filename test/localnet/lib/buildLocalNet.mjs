import pEvent from "p-event";
import { $, sleep } from "zx";
import { initAllChains } from "./initAllChains.mjs";
import { initAllRelayers } from "./initAllRelayers.mjs";
import { startAllChains } from "./startAllChains.mjs";
import { startAllRelayers } from "./startAllRelayers.mjs";
import { takeSnapshot } from "./takeSnapshot.mjs";

export async function buildLocalNet({
  network,
  home = "/tmp/localnet",
  registryFrom = `/tmp/localnet/registry`,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
}) {
  // 1) init all IBC chains
  await initAllChains({ network, home });
  // 2) start all IBC chains
  const chainsProps = await startAllChains({
    network,
    home,
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
  });
  // 2b) wait for IBC chains first block to be written
  await Promise.all(
    chainsProps.map(async ({ proc }) => {
      const asyncIterator = pEvent.iterator(proc.stderr, "data", {
        resolutionEvents: ["finish"],
      });
      for await (let chunk of asyncIterator) {
        if (chunk.includes("indexed block")) break;
      }
    })
  );
  await sleep(5000);
  await $`curl http://localhost:11000`;
  // 3) init all IBC relayers
  await initAllRelayers({ network, home, registryFrom });
  // 4) start all IBC relayers
  const relayersProps = await startAllRelayers({
    network,
    home,
    rpcInitialPort,
    p2pInitialPort,
    pprofInitialPort,
  });
  // 5) wait for IBC relayers confirmation message and then stop relayers
  await Promise.all(
    relayersProps.map(async ({ proc }) => {
      const asyncIterator = pEvent.iterator(proc.stdout, "data", {
        resolutionEvents: ["finish"],
      });
      for await (let chunk of asyncIterator) {
        if (chunk.includes("next heights to relay")) break;
      }
      proc.kill("SIGTERM");
    })
  );
  // 6) now stop IBC chains
  await Promise.all(
    chainsProps.map(async ({ proc }) => {
      proc.kill("SIGTERM");
    })
  );
  await takeSnapshot({ home });
}
