import { LcdClient } from "@cosmjs/launchpad";

export interface ClpExtension {
  readonly clp: {};
}

export function setupClpExtension(base: LcdClient): ClpExtension {
  return {
    clp: {},
  };
}
