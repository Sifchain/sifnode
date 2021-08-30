import { SynchronousCommand, SynchronousCommandResult } from "./synchronousCommand";
import { DeployedContractAddresses } from "../../scripts/deploy_contracts";

export class SmartContractDeployResult extends SynchronousCommandResult {
  constructor(
    readonly contractAddresses: DeployedContractAddresses,
    readonly completed: boolean,
    readonly error: Error | undefined,
    readonly output: string
  ) {
    super(completed, error, output);
  }
}

export class SmartContractDeployResultsPromise {
  constructor(
    readonly results: Promise<SmartContractDeployResult>
  ) {
  }
}

export class SmartContractDeployer extends SynchronousCommand<SmartContractDeployResult> {
  constructor() {
    super();
  }

  cmd(): [string, string[]] {
    return ["npx", [
      "hardhat",
      "run",
      "scripts/deploy_contracts.ts",
      "--network",
      "localhost"
    ]]
  }

  resultConverter(r: SynchronousCommandResult): SmartContractDeployResult {
    // This is to handle npx commmand outputting "No need to generate any newer types"
    const jsonOutput = JSON.parse(r.output.split('\n')[1]);
    return new SmartContractDeployResult({
      bridgeBank: jsonOutput.bridgeBank,
      bridgeRegistry: jsonOutput.bridgeRegistry,
      rowanContract: jsonOutput.rowanContract
    },
      r.completed, r.error, r.output);
  }
}
