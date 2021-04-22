#!/usr/bin/env node

import { spawn, exec } from "child_process";
import { resolve } from "path";
import { sleep } from "../e2e/utils";
import treekill from "tree-kill";
import chalk from "chalk";

const cmdStack = [];
/**
 * This is a utility to be used within our testing frameworks for restarting our backend stack.
 */
export async function restartStack() {
  const uiFolder = resolve(__dirname, "../");

  const cmd = spawn("./scripts/run-stack-backend.sh", [], { cwd: uiFolder });
  cmdStack.push(cmd);
  let handler;

  await new Promise<void>((resolve) => {
    handler = (data) => {
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

function treeKillProm(pid) {
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
