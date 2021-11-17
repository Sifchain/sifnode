import { $, cd } from "zx";

export async function createArchive({ src, basePath = "/tmp", name }) {
  cd(`/`);
  await $`rm -f ${basePath}/${name}.tbz`;
  await $`tar cvjf ${basePath}/${name}.tbz ${src}`;
}
