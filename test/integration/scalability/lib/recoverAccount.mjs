import { $, nothrow } from "zx";

import { createRequire } from "module";
const require = createRequire(import.meta.url);

const chains = require("../config/chains.json");
const accounts = require("../config/accounts.json");

export async function recoverAccount({ name }) {
  if (!name) throw new Error("missing requirement argument: --name");
  if (!accounts[name])
    throw new Error("account name does not exist in accounts.json");

  const mnemonic = accounts[name];
  const chainsArr = Object.values(chains).filter(
    ({ disabled = false }) => !disabled
  );

  for (let i = 0; i < chainsArr.length; i++) {
    const { binary } = chainsArr[i];
    await nothrow(
      $`${binary} keys delete ${name} --keyring-backend test -y 2> /dev/null`
    );
    await $`printf "%s\\n\\n" ${mnemonic} | ${binary} keys add ${name} -i --recover --keyring-backend test`;
  }
}
