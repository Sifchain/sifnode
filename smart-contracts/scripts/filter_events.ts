import { ethers } from "hardhat";
import { promises } from "fs"
import { BigNumber } from "ethers";
import { CoinGeckoClient } from "coingecko-api-v3";
import { Client } from "pg";


const BridgeBankAddress = "0xB5F54ac4466f5ce7E0d8A5cB9FE7b8c0F35B7Ba8";
const AttackerCosmosAddresses = ["sif1vmzzvtwr5dl2dumh0l65fsvfdf6j2etv2tt4xy", "sif164lcv2wxzyzy6g4ea7c2jrwjhqukfll8jr6wa2"];
const EthereumStartingHeight = 14643634;

function bytesToSifAddress(bytes: string) : string {
    const address = ethers.utils.toUtf8String(bytes);
    // console.log("Parsed Cosmos Address: ", address);
    return address;
}

interface BridgeBankEvent {
    BlockNumber: number;
    TransactionHash: string;
    EthereumAddress: string;
    TokenSymbol: string;
    TokenAddress: string;
    Value: string;
    FormattedValue: string;
}

interface BurnOrLockEvent extends BridgeBankEvent {
    CosmosAddress: string;
    Nonce: BigNumber;
}

type UnlockOrMintEvent = BridgeBankEvent;

const client = new CoinGeckoClient({
  timeout: 10000,
  autoRetry: true,
});

// USERNAME PASSWORD HOST PORT DATABASE
const pgClient = new Client();

async function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms))
}

function ParseValue(symbol: string, value: BigNumber) : string {
    let decimals = 18;
    switch (symbol) {
        case "USDC":
            decimals = 6;
            break;
        default: 
            decimals = 18;
            break;
    }
    const formatted = ethers.utils.formatUnits(value, decimals);
    return ethers.utils.commify(formatted);
}

async function main () {

    const bridgeBankFactory = await ethers.getContractFactory("BridgeBank");
    const cosmosBridgeFactory = await ethers.getContractFactory("CosmosBridge");
    const bridgeBank = await bridgeBankFactory.attach(BridgeBankAddress);
    const cosmosBridgeAddress = await bridgeBank.cosmosBridge();
    const cosmosBridge = await cosmosBridgeFactory.attach(cosmosBridgeAddress);
    console.log("Connected to CosmosBridge at address: ", cosmosBridgeAddress);

    const BurnEvents = await bridgeBank.queryFilter(bridgeBank.filters.LogBurn(), EthereumStartingHeight);
    const LockEvents = await bridgeBank.queryFilter(bridgeBank.filters.LogLock(), EthereumStartingHeight);
    const UnlockEvents = await bridgeBank.queryFilter(bridgeBank.filters.LogUnlock(), EthereumStartingHeight);
    const MintEvents = await bridgeBank.queryFilter(bridgeBank.filters.LogBridgeTokenMint(), EthereumStartingHeight);
    
    async function ParseLockBurn(BurnOrLock: typeof BurnEvents[0] | typeof LockEvents[0]): Promise<BurnOrLockEvent> {
        const LockOrBurn = {
            BlockNumber: BurnOrLock.blockNumber,
            TransactionHash: BurnOrLock.blockHash,
            CosmosAddress: bytesToSifAddress(BurnOrLock.args?._to),
            EthereumAddress: BurnOrLock.args?._from,
            TokenAddress: BurnOrLock.args?._token,
            TokenSymbol: BurnOrLock.args?._symbol,
            Value: BurnOrLock.args?._value.toString(),
            FormattedValue: ParseValue(BurnOrLock.args?._symbol, BurnOrLock.args?._value),
            Nonce: BurnOrLock.args?._nonce
        };
        return LockOrBurn;
    }

    async function ParseUnlock(Unlock: typeof UnlockEvents[0]): Promise<UnlockOrMintEvent> {
        const UnlockEvent = {
            BlockNumber: Unlock.blockNumber,
            TransactionHash: Unlock.blockHash,
            EthereumAddress: Unlock.args?._to,
            TokenAddress: Unlock.args?._token,
            TokenSymbol: Unlock.args?._symbol,
            Value: Unlock.args?._value.toString(),
            FormattedValue: ParseValue(Unlock.args?._symbol, Unlock.args?._value),
        };
        return UnlockEvent;
    }

    async function ParseMint(Mint: typeof MintEvents[0]): Promise<UnlockOrMintEvent> {
        const MintEvent = {
            BlockNumber: Mint.blockNumber,
            TransactionHash: Mint.transactionHash,
            EthereumAddress: Mint.args?._beneficiary,
            TokenAddress: Mint.args?._token,
            TokenSymbol: Mint.args?._symbol,
            Value: Mint.args?._amount.toString(),
            FormattedValue: ParseValue(Mint.args?._symbol, Mint.args?._amount),
        };
        return MintEvent;
    }

    const Burns = await Promise.all(BurnEvents.map(async (burn) => await ParseLockBurn(burn)));
    const Locks = await Promise.all(LockEvents.map(async (lock) => await ParseLockBurn(lock)));
    const Unlocks = await Promise.all(UnlockEvents.map(async (unlock) => await ParseUnlock(unlock)));
    const Mints = await Promise.all(MintEvents.map(async (mint) => await ParseMint(mint)));
    const Exports = await cosmosBridge.queryFilter(cosmosBridge.filters.LogProphecyCompleted(), EthereumStartingHeight);
    const AttackerBurns : BurnOrLockEvent[] = [];
    const AttackerLocks : BurnOrLockEvent[] = [];
    const AttackerUnlocks : UnlockOrMintEvent[] = [];
    const AttackerMints : UnlockOrMintEvent[] = [];
    
    for (const AttackerCosmosAddress of AttackerCosmosAddresses) {
        AttackerBurns.push(...Burns.filter((burn) => burn.CosmosAddress === AttackerCosmosAddress));
        AttackerLocks.push(...Locks.filter((lock) => lock.CosmosAddress === AttackerCosmosAddress));
    }
    
    const AttackerEthereumAddresses = Array.from(new Set<string>([...AttackerBurns.map(burn => burn.EthereumAddress), ...AttackerLocks.map(lock => lock.EthereumAddress)]));
    for (const ethereumAddress of AttackerEthereumAddresses) {
        AttackerUnlocks.push(...Unlocks.filter((unlock) => unlock.EthereumAddress === ethereumAddress));
        AttackerMints.push(...Mints.filter((mint) => mint.EthereumAddress === ethereumAddress));
    }

    const output = {
        AttackerCosmosAddresses,
        AttackerEthereumAddresses,
        Burns,
        Locks,
        Unlocks,
        Mints,
        Exports,
        AttackerBurns,
        AttackerLocks,
        AttackerUnlocks,
        AttackerMints
    };

    const jsonOutput = JSON.stringify(output);
    await promises.writeFile("imports_exports.json", jsonOutput);
}

main().then(() => console.log("Exited successfully")).catch((error) => console.error("Error Reported: ", error))