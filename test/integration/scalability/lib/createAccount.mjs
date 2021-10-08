import { recoverAccount } from "./recoverAccount.mjs";

const accounts = require("./accounts.json");

export async function createAccount({ name }) {
  if (!name) throw new Error("missing requirement argument: --name");
  if (accounts[name])
    throw new Error("account name already exists in accounts.json");

  const { stdout: mnemonic } = await nothrow(
    $`yes | sifnoded keys add ${name} --keyring-backend test --output json | jq -r .mnemonic`
  );
  accounts[name] = mnemonic.replace("\n", "");
  await fs.writeFile("./accounts.json", JSON.stringify(accounts, null, 2));

  await recoverAccount({ name });
}
