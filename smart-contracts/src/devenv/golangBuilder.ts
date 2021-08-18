import {registry, singleton} from "tsyringe";
import * as childProcess from "child_process"
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import * as hre from "hardhat"
import {EthereumAccounts, EthereumResults, ShellCommand} from "./devEnv";

interface GolangArguments {
}

@singleton()
export class GolangBuilder extends ShellCommand {
    constructor(
        readonly args: GolangArguments
    ) {
        super();
    }

    cmd(): [string, string[]] {
        return ["make", [
            "install",
        ]]
    }

    override run(): Promise<void> {
        return new Promise((resolve, reject) => {
            const [c, args] = this.cmd()
            const process = childProcess.spawn(c, args, {stdio: "inherit"})
            process.on("error", err => {
                reject(err)
            })
            process.on("exit", e => {
                if (e && e != 0) {
                    reject(e)
                } else
                    resolve()
            })
        })
    }
}
