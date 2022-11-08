import pEvent from "p-event";
import { $, sleep } from "zx";
import { initAllChains } from "./initAllChains.mjs";
import { initAllRelayers } from "./initAllRelayers.mjs";
import { startAllChains } from "./startAllChains.mjs";
import { startAllRelayers } from "./startAllRelayers.mjs";
import { takeSnapshot } from "./takeSnapshot.mjs";

export async function buildLocalNet({
  network,
  configPath = "/tmp/localnet/config",
  registryFrom = `/tmp/localnet/config/registry`,
  rpcInitialPort = 11000,
  p2pInitialPort = 12000,
  pprofInitialPort = 13000,
}) {
  $.verbose = false;
  // 1) init all IBC chains
  await initAllChains({ network, configPath });
  // 2) start all IBC chains
  const chainsProps = await startAllChains({
    network,
    configPath,
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
  $.verbose = true;
  // 3) init all IBC relayers
  await initAllRelayers({ network, configPath, registryFrom });
  // 4) start all IBC relayers
  const relayersProps = await startAllRelayers({
    network,
    configPath,
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
      proc.kill();
    })
  );
  // 6) now stop IBC chains
  await Promise.all(
    chainsProps.map(async ({ proc }) => {
      proc.kill();
    })
  );
  await takeSnapshot({ configPath });
}
