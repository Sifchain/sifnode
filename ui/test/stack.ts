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
export async function restartStack() {
  const uiFolder = resolve(__dirname, "../");

  const cmd = spawn("./scripts/run-stack-backend.sh", [], { cwd: uiFolder });
  cmdStack.push(cmd);
  let handler;

  await new Promise<void>((resolve) => {
    handler = (data: Buffer) => {
      const dataStr = data.toString().replace(/\n$/, "");
      console.log(chalk.blue(dataStr));
      if (dataStr.includes("cosmos process events for blocks")) {
        resolve();
      }
    };
    cmd.stdout.on("data", handler);
    cmd.stderr.on("data", handler);
    cmd.on("error", (err) => {
      console.log(chalk.red(err.toString()));
    });
  });
  console.log(chalk.green("DONE"));
}

function treeKillProm(pid: number) {
  return new Promise((resolve) => {
    treekill(pid, resolve);
  });
}

export async function killStack() {
  await Promise.all(cmdStack.map((cmd) => treeKillProm(cmd.pid)));

  await new Promise((resolve) => {
    exec("killall sifnoded sifnodecli ebrelayer ganache-cli", resolve);
    console.log("⬇⬇⬇  S.T.A.C.K  ⬇⬇⬇");
  });
  await sleep(1000);
}

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
    afterEach(async () => await restartStack());
  }
  afterAll(async () => await killStack());
}
