#!/usr/bin/env zx

import { takeSnapshot } from "../lib/takeSnapshot.mjs";
import { arg } from "../utils/arg.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

const args = arg(
  {
    "--home": String,
  },
  `
Usage:

  yarn takeSnapshot [options]

Create a snapshot of all the localnet file-based data including the IBC chains + relayers.

Options:

--home      Global directory for config and data of initiated chains
`
);

const home = args["--home"] || undefined;

const chainProps = getChainProps({
  home,
});

await takeSnapshot({ ...chainProps });
