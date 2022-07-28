#!/usr/bin/env zx

$.verbose = false;

import { getEntries } from "./helpers/getEntries.mjs";
import { getFunds } from "./helpers/getFunds.mjs";
import { getTokens } from "./helpers/getTokens.mjs";

const entries = await getEntries();
const tokens = await getTokens(entries);
const funds = getFunds(tokens, "\\,");

echo(funds);
