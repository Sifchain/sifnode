#!/usr/bin/env zx

import { recoverAccount } from "../lib/recoverAccount.mjs";
import { arg } from "../utils/arg.mjs";

const args = arg(
  {
    "--name": String,
  },
  `
Usage:

  yarn recoverAccount [options]

Recover one account.

Options:

--name      Account name or address
`
);

const name = args["--name"] || undefined;

await recoverAccount({ name });
