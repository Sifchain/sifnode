import { spinner } from "zx/experimental";

export async function getEntries() {
  const {
    result: {
      registry: { entries },
    },
  } = await spinner("loading entries                              ", () =>
    within(
      async () =>
        // (await fetch("https://api.sifchain.finance/tokenregistry/entries")).json()
        await fs.readJson(`./entries.json`)
    )
  );
  return entries;
}
