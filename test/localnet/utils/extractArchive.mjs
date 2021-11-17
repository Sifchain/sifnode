import { $, cd } from "zx";

export async function extractArchive({ dst = "/", basePath = "/tmp", name }) {
  cd(`${dst}`);
  await $`tar xvjf ${basePath}/${name}.tbz`;
}
