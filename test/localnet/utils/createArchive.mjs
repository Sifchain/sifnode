import { $, cd } from "zx";

export async function createArchive({
  configPath = "/tmp/localnet/config",
  archivePath = "/tmp/localnet/config.tbz",
}) {
  cd(configPath);
  await $`rm -f ${archivePath}`;
  await $`tar -cjf ${archivePath} .`;
}
