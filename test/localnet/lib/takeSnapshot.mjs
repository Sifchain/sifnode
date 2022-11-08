import { createArchive } from "../utils/createArchive.mjs";

export async function takeSnapshot({
  configPath = `/tmp/localnet/config`,
  archivePath = `/tmp/localnet/config.tbz`,
}) {
  await createArchive({ configPath, archivePath });
}
