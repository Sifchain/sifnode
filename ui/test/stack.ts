#!/usr/bin/env node

import { spawn, ChildProcess, exec } from "child_process";
import { resolve } from "path";

import treekill from "tree-kill";
import chalk from "chalk";
const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const cmdStack: ChildProcess[] = [];
/**
 * This is a utility to be used within our testing frameworks for restarting our backend stack.
 */
const uiFolder = resolve(__dirname, "../");
export async function restartStack() {
  if (process.env.NOSTACK) return;
  console.log("^^^  START RESTART  ^^^");

  const cmd = spawn("./scripts/run-stack-backend.sh", [], {
    cwd: uiFolder,
    stdio: ["ignore", "pipe", "pipe"],
    detached: true,
  });
  cmdStack.push(cmd);
  let handler: (d: Buffer) => void;

  await new Promise<void>((resolve) => {
    const errHandler = (err: Error) => {
      console.log(chalk.red(err.toString()));
    };
    handler = (data: Buffer) => {
      const dataStr = data.toString().replace(/\n$/, "");
      console.log(chalk.blue(dataStr));
      if (dataStr.includes("cosmos process events for blocks")) {
        cmd.stdout.off("data", handler);
        cmd.stderr.off("data", handler);
        cmd.off("error", errHandler);
        setTimeout(resolve, 3000); // Adding a timeout here seems to stabilize keplr
      }
    };
    cmd.stdout.on("data", handler);
    cmd.stderr.on("data", handler);
    cmd.on("error", errHandler);
  });
  console.log(chalk.green("DONE"));
}

function treeKillProm(pid: number) {
  return new Promise((resolve) => {
    treekill(pid, resolve);
  });
}

export async function killStack() {
  if (process.env.NOSTACK) return;
  console.log("⬇⬇⬇  START SHUTDOWN  ⬇⬇⬇");
  console.log(cmdStack.map((c) => c.pid).join(":"));

  await new Promise<void>((resolve, reject) => {
    exec("./scripts/kill-stack-backend.sh", { cwd: uiFolder }, (err, out) => {
      if (err) return reject(err);
      resolve();
    });
  });

  for (let cmd of cmdStack) {
    console.log("...killing");
    await treeKillProm(cmd.pid);
    console.log("...finished killing");
  }

  console.log("⬇⬇⬇  S.T.A.C.K  ⬇⬇⬇");
}

/**
 * Utility for tooling tests with our backing stack. When `when` is `once` the stack is
 * booted up at the start of the file and then torn down. When 'every-test` is passed in then a new stack gets created for every test.
 * @param when "once" or "every-test"
 */
export function useStack(when: "once" | "every-test") {
  // This might change if we work out a way to run each jest test in it's own container.
  // For now this is just a sanity check as it is easy to accidentally mess up here.
  if (
    process.argv.filter((arg) => ["--runInBand", "-i"].includes(arg)).length ===
    0
  ) {
    throw new Error(
      "To use the sifchain stack in a test you must run them in band  with either the `-i` or `--runInBand` flags eg.: \n\n\tjest -i MyTest.test",
    );
  }

  beforeAll(async () => await restartStack());
  if (when === "every-test") {
    afterEach(async () => {
      await restartStack();
      await sleep(1000);
    });
  }
  afterAll(async () => await killStack());
}
