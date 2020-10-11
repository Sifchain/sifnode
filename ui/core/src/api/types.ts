import * as io from "./index";

export type FullApi = typeof io;

export type Api<
  T extends keyof FullApi = keyof FullApi,
  U extends object = {}
> = {
  api: Pick<FullApi, T>;
} & U;
