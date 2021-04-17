#!/usr/bin/env node

import { spawn } from "child_process";
import { resolve } from "path";
import kill from "tree-kill";

export async function restartStack(verbose = false) {
  console.log("...launching stack");
  const uiFolder = resolve(__dirname, "../../");
  const cmd = spawn("./scripts/run-stack-backend.sh", [], { cwd: uiFolder });

  await new Promise<void>((resolve) => {
    const handler = (data: Buffer) => {
      const dataStr = data.toString().replace(/\n$/, "");
      verbose && console.log(dataStr);
      if (dataStr.includes("cosmos process events for blocks")) {
        console.log("⬆⬆⬆  S.T.A.C.K  ⬆⬆⬆");
        !verbose && cmd.stdout.off("data", handler);
        !verbose && cmd.stderr.off("data", handler);
        resolve();
      }
    };
    cmd.stdout.on("data", handler);
    verbose && cmd.stderr.on("data", handler);
  });

  return async function killStack() {
    if (!cmd.killed) {
      await new Promise((done) => kill(cmd.pid, done));
      console.log("⬇⬇⬇  S.T.A.C.K  ⬇⬇⬇");
    }
  };
}

export function withStack(handler: () => Promise<any>, verbose = false) {
  return async () => {
    const VERBOSE = !!process.env.VERBOSE;
    const kill = await restartStack(verbose || VERBOSE);
    await handler();
    await kill();
  };
}
