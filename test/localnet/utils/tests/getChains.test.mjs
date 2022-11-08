import { getChains } from "../getChains.mjs";

import { createRequire } from "module";
const require = createRequire(import.meta.url);
const defaultChains = require("../../config/chains.json");

test("gets chains", () => {
  const chains = getChains({ chains: defaultChains });

  expect(chains).toMatchSnapshot();
});
