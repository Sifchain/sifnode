import {singleton} from "tsyringe";
import * as childProcess from "child_process"
import {GolangResults, ShellCommand} from "./devEnv";
import {firstValueFrom, ReplaySubject} from "rxjs";

@singleton()
class GolangArguments {
}

// export function resolvableAndRejectablePromise<T>(): [Promise<T>, (value: (PromiseLike<T> | T)) => void, (value: any) => void] {
//     let resolveFn: (value: (PromiseLike<T> | T)) => void
//     let rejectFn: (value: any) => void
//     const result = new Promise<T>((resolve, reject) => {
//         resolveFn = resolve
//         rejectFn = reject
//     })
//     return [result, resolveFn, rejectFn]
// }
//
// export function resolvablePromise<T>(): [Promise<T>, (T) => void] {
//     const [promise, resolveFn, _] = resolvableAndRejectablePromise<T>()
//     return [promise, resolveFn]
// }


@singleton()
export class GolangBuilder extends ShellCommand {
    constructor(
        readonly args: GolangArguments
    ) {
        super();
    }

    cmd(): [string, string[]] {
        return ["make", [
            "-C",
            "..",
            "install",
        ]]
    }

    completion = new ReplaySubject<GolangResults>(1)

    override async run(): Promise<void> {
        const pq = childProcess.spawnSync(this.cmd()[0], this.cmd()[1])
        this.completion.next({golangBuilt: true, output: pq.stdout.toString()})
        return Promise.resolve()
    }

    results(): Promise<GolangResults> {
        return firstValueFrom(this.completion)
    }
}
