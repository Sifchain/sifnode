import { $ } from "zx";

export async function getAddress({
  binary,
  name,
  home = undefined,
  binPath = "/tmp/localnet/bin",
}) {
  if (!binary) throw new Error("missing requirement argument: --binary");
  if (!name) throw new Error("missing requirement argument: --name");
  if (!binPath) throw new Error("missing requirement argument: --binPath");

  const addr = await $`
${binPath}/${binary} \
    keys \
    show \
    ${name} \
    --keyring-backend test \
    -a \
    ${home ? `--home` : ``} ${home ? `${home}` : ``} \
    2> /dev/null || echo ${name}`;
  return addr;
}
