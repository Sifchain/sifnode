import { getChains } from "../getChains.mjs";
const defaultChains = require("../../config/chains.json");

function test() {
  const chains = getChains({ chains: defaultChains });

  console.log(chains);
}

test();
