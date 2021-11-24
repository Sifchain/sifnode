import { $, cd } from "zx";

export async function extractArchive({
  configPath = "/tmp/localnet/config",
  archivePath = "/tmp/localnet/config.tbz",
}) {
  await $`mkdir -p ${configPath}`;
  cd(`${configPath}`);
  await $`tar -xjf ${archivePath}`;
}
