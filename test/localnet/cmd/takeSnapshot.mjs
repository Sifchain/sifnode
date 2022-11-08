import { takeSnapshot } from "../lib/takeSnapshot.mjs";
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

  yarn takeSnapshot [options]

Create a snapshot of all the localnet file-based data including the IBC chains + relayers.

Options:

--home      Global directory for config and data of initiated chains
--configPath       Location of the config directory where chains and relayers data are stored
--archivePath      Location where the snapshot archive file will be created
`
  );

  const configPath = args["--configPath"] || undefined;
  const archivePath = args["--archivePath"] || undefined;

  const chainProps = getChainProps({
    configPath,
    archivePath,
  });

  await takeSnapshot({ ...chainProps });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
