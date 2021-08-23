import {HardhatRuntimeEnvironment} from "hardhat/types";
import {HardhatRuntimeEnvironmentToken,} from "./injectionTokens";
import {SignerWithAddress} from "@nomiclabs/hardhat-ethers/signers";
import {inject, injectable} from "tsyringe";
import {isHardhatRuntimeEnvironment} from "./hardhatSupport";

/**
 * The accounts necessary for testing a sifchain system
 */
export class SifchainAccounts {
    constructor(
        readonly operatorAccount: SignerWithAddress,
        readonly ownerAccount: SignerWithAddress,
        readonly pauserAccount: SignerWithAddress,
        readonly validatatorAccounts: Array<SignerWithAddress>,
        readonly availableAccounts: Array<SignerWithAddress>
    ) {
    }
}

/**
 * Note that the hardhat environment provides accounts as promises, so
 * we need to wrap a SifchainAccounts in a promise.
 */
@injectable()
export class SifchainAccountsPromise {
    accounts: Promise<SifchainAccounts>

    constructor(accounts: Promise<SifchainAccounts>);
    constructor(@inject(HardhatRuntimeEnvironmentToken) hardhatOrAccounts: HardhatRuntimeEnvironment | Promise<SifchainAccounts>) {
        if (isHardhatRuntimeEnvironment(hardhatOrAccounts)) {
            this.accounts = hreToSifchainAccountsAsync(hardhatOrAccounts)
        } else {
            this.accounts = hardhatOrAccounts
        }
    }
}

export async function hreToSifchainAccountsAsync(hardhat: HardhatRuntimeEnvironment): Promise<SifchainAccounts> {
    const accounts = await hardhat.ethers.getSigners()
    const [operatorAccount, ownerAccount, pauserAccount, validator1Account, ...extraAccounts] = accounts
    return new SifchainAccounts(operatorAccount, ownerAccount, pauserAccount, [validator1Account], extraAccounts)
}
