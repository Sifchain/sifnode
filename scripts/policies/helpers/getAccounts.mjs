import { spinner } from "zx/experimental";

export async function getAccounts() {
  const accounts = await spinner(
    "loading accounts                              ",
    () => within(async () => await fs.readJson(`./accounts.json`))
  );
  return accounts;
}
