import { downloadBinaries } from "../lib/downloadBinaries.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

export async function start() {
  const args = arg(
    {
      "--home": String,
    },
    `
Usage:

  yarn downloadBinaries [options]

Download all the binaries of IBC chains.

Options:

--home      Global directory for config and data of initiated chains
`
  );

  const home = args["--home"] || undefined;

  const chainProps = getChainProps({
    home,
  });

  await downloadBinaries({ ...chainProps });
}

if (process.env.NODE_ENV !== "test") {
  start();
}
