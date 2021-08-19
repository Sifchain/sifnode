import {singleton} from "tsyringe";
import {SynchronousCommand, SynchronousCommandResult} from "./synchronousCommand";

export class GolangResults extends SynchronousCommandResult {
}

@singleton()
export class GolangBuilder extends SynchronousCommand<GolangResults> {
    constructor(
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

    resultConverter = (r: SynchronousCommandResult) => new GolangResults(r.completed, r.error, r.output)
}
