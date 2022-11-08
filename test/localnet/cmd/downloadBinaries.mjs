import { downloadBinaries } from "../lib/downloadBinaries.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--binPath": String,
    },
    `
Usage:

  yarn downloadBinaries [options]

Download all the binaries of IBC chains.

Options:

--binPath      Global directory for binaries location
`
  );

  const binPath = args["--binPath"] || undefined;

  const chainProps = getChainProps({
    binPath,
  });

  await downloadBinaries({ ...chainProps });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
