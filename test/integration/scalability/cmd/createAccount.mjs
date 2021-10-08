#!/usr/bin/env zx

import { createAccount } from "../lib/createAccount.mjs";
import { arg } from "../utils/arg.mjs";

const args = arg(
  {
    "--name": String,
  },
  `
Usage:

  yarn createAccount [options]

Create a new account.

Options:

--name      Account name or address
`
);

const name = args["--name"] || undefined;

await createAccount({ name });
