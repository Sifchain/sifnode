import { SynchronousCommand, SynchronousCommandResult } from "./synchronousCommand"
import { requiredEnvVar } from "../contractSupport"

export class GolangResults extends SynchronousCommandResult {
  constructor(
    readonly goBin: string,
    readonly completed: boolean,
    readonly error: Error | undefined,
    readonly output: string
  ) {
    super(completed, error, output)
  }
}

export class GolangResultsPromise {
  constructor(readonly results: Promise<GolangResults>) {}
}

export class GolangBuilder extends SynchronousCommand<GolangResults> {
  constructor() {
    super()
  }

  cmd(): [string, string[]] {
    return ["make", ["-C", "..", "install"]]
  }

  resultConverter(r: SynchronousCommandResult): GolangResults {
    const goBin = requiredEnvVar("GOBIN")
    return new GolangResults(goBin, r.completed, r.error, r.output)
  }
}
