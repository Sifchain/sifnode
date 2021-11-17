import { loadLocalNet } from "../lib/loadLocalNet.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--basePath": String,
      "--name": String,
    },
    `
Usage:

  yarn loadLocalnet [options]

Load an existing snapshot of localnet to run chains + relayers on top.

Options:

--basePath  Location of the snapshot archive
--name      Name of the snapshot to load
`
  );

  const basePath = args["--basePath"] || undefined;
  const name = args["--name"] || undefined;

  const chainProps = getChainProps({
    basePath,
    name,
  });

  await loadLocalNet({ ...chainProps });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
