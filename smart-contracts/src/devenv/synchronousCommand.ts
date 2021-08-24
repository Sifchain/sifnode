import {singleton} from "tsyringe";
import * as ChildProcess from "child_process"
import {ShellCommand} from "./devEnv";
import {firstValueFrom, ReplaySubject} from "rxjs";

export class SynchronousCommandResult {
    constructor(
        readonly completed: boolean,
        readonly error: Error | undefined,
        readonly output: string
    ) {

    }
}

export abstract class SynchronousCommand<T extends SynchronousCommandResult> extends ShellCommand<T> {
    protected constructor() {
        super();
    }

    completion = new ReplaySubject<T>(1)

    abstract resultConverter(x: SynchronousCommandResult): T

    override async run(): Promise<void> {
        const commandResult = ChildProcess.spawnSync(this.cmd()[0], this.cmd()[1])
        let synchronousCommandResult = new SynchronousCommandResult(
            true,
            commandResult.error,
            commandResult.stdout?.toString() ?? ""
        );
        this.completion.next(this.resultConverter(synchronousCommandResult))
        return Promise.resolve()
    }

    results(): Promise<T> {
        return firstValueFrom(this.completion)
    }
}

export async function sampleCode() {
    // tldr: execSync for short commands without much output.  spawn for anything long-running.

    // What the doc says:

    // spawn is the base function.
    // exec passes stdin and stdout to a callback when the command completes; with async you
    //   probably never want this
    // *File* skips spawning a shell
    // fork is a node-specific thing
    //
    // What you actually care about
    //
    // spawnSync returns an object with fields like status and stderr.  Use spawn with stdio: "inherit"
    //   for long-lived processes
    //
    // execSync just returns a string; use this for simple, short commands.  The doc is misleading
    //   about defaulting to encoding: utf8; that's not really what it does.  Specify encoding.
    //
    // execSync raises an exception; spawnSync sets error fields on the result object
    {
        // A short, synchronous command with input and output as simple variables:
        // const result = ChildProcess.execSync("wc -c", {
        //     encoding: "ascii",
        //     input: "abc"
        // })
        // console.log("result length for wc -c of abc should be 3: ", result.trim())
        //
        // // result1 returns an object with a status of 126, fields stdout and stderr set.
        // // Use execSync if you just want to get back a string, and you're ok with an exception
        // const result1 = ChildProcess.spawnSync("fnord ls -l /tmp", {
        //     shell: true,
        //     encoding: "ascii"
        // })
        //
        // const result4 = await ChildProcess.spawn("ls -l /tmp", {
        //     stdio: "inherit"
        // })
        //
        // // execSync with a bad command throws an exception
        // try {
        //     const result2 = await ChildProcess.execSync("ls -l /tmp", {encoding: "ascii"})
        // } catch (e) {
        //     console.log("got an exception")
        // }
        //
        // const result3 = await ChildProcess.spawnSync("wc", {
        //     encoding: "ascii",
        //     shell: true,
        //     input: "abc"
        // })
        // console.log("result3 should be the result of echo abc | wc", result3)
        //
        // console.log("result1: ", result1)
    }
}
