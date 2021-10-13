import { initAllChains } from "./initAllChains.mjs";
import { takeSnapshot } from "./takeSnapshot.mjs";

export async function buildLocalnet({ network, home = "/tmp/localnet" }) {
  await initAllChains({ network, home });
  await takeSnapshot({ home });
}
