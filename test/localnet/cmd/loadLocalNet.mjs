import { loadLocalNet } from "../lib/loadLocalNet.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--configPath": String,
      "--archivePath": String,
    },
    `
Usage:

  yarn loadLocalnet [options]

Load an existing snapshot of localnet to run chains + relayers on top.

Options:

--configPath       Location of the config directory where chains and relayers data will be stored
--archivePath      Location of the snapshot archive file that contains all the config files
`
  );

  const configPath = args["--configPath"] || undefined;
  const archivePath = args["--archivePath"] || undefined;

  const chainProps = getChainProps({
    configPath,
    archivePath,
  });

  await loadLocalNet({ ...chainProps });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
