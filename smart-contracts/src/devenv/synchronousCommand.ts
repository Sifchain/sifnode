import * as ChildProcess from "child_process"
import { ShellCommand } from "./devEnv"
import { firstValueFrom, ReplaySubject } from "rxjs"

export class SynchronousCommandResult {
  constructor(
    readonly completed: boolean,
    readonly error: Error | undefined,
    readonly output: string
  ) {}
}

export abstract class SynchronousCommand<
  T extends SynchronousCommandResult
> extends ShellCommand<T> {
  protected constructor() {
    super()
  }

  completion = new ReplaySubject<T>(1)

  abstract resultConverter(x: SynchronousCommandResult): T

  override async run(): Promise<void> {
    const commandResult = ChildProcess.spawnSync(this.cmd()[0], this.cmd()[1])
    let synchronousCommandResult = new SynchronousCommandResult(
      true,
      commandResult.error,
      commandResult.stdout?.toString() ?? ""
    )
    this.completion.next(this.resultConverter(synchronousCommandResult))
    return Promise.resolve()
  }

  results(): Promise<T> {
    return firstValueFrom(this.completion)
  }
}
