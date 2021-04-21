#!/usr/bin/env node

const { spawn, exec } = require("child_process");
const { resolve } = require("path");
const treekill = require("tree-kill");
const chalk = require("chalk");
const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));
const cmdStack = [];
async function restartStack() {
  const uiFolder = resolve(__dirname, "../");

  const cmd = spawn("./scripts/run-stack-backend.sh", [], { cwd: uiFolder });
  cmdStack.push(cmd);
  let handler;

  await new Promise((resolve) => {
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
      console.log(chalk.red(err));
    });
  });
  console.log(chalk.green("DONE"));
}

function treeKillProm(pid) {
  return new Promise((resolve) => {
    treekill(pid, resolve);
  });
}

async function killStack() {
  await Promise.all(cmdStack.map((cmd) => treeKillProm(cmd.pid)));

  await new Promise((resolve) => {
    exec(
      "killall",
      ["sifnoded", "sifnodecli", "ebrelayer", "ganache-cli"],
      (out) => {
        resolve(out);
      },
    );
    console.log("⬇⬇⬇  S.T.A.C.K  ⬇⬇⬇");
  });
  await sleep(1000);
}

module.exports.restartStack = restartStack;

module.exports.killStack = killStack;

// restartStack().then(() => {
//   console.log("Next");
// });
