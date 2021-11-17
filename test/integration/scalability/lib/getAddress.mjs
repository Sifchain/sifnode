import { $ } from "zx";

export async function getAddress({ binary, name, home = undefined }) {
  if (!binary) throw new Error("missing requirement argument: --binary");
  if (!name) throw new Error("missing requirement argument: --name");

  const addr = await $`
${binary} \
    keys \
    show \
    ${name} \
    --keyring-backend test \
    -a \
    ${home ? `--home` : ``} ${home ? `${home}` : ``} \
    2> /dev/null || echo ${name}`;
  return addr;
}
