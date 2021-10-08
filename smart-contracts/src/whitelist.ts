import {BridgeBank, BridgeToken__factory} from "../build";

interface LogWhiteListUpdateRawData {
    token: string
    enabled: boolean,
    blockNumber: number
    transactionIndex: number
}

interface Erc20Data {
    name: string
    symbol: string
    decimals: number
    address: string
}

interface WhitelistItem {
    logData: LogWhiteListUpdateRawData
    erc20Data: Erc20Data
}

/**
 * Get the raw log entries for whitelisted tokens.
 *
 * Use getWhitelistItems if you also need ERC20 data.
 *
 * @param bridgeBank
 */
export async function getWhitelistRawItems(
    bridgeBank: BridgeBank
): Promise<LogWhiteListUpdateRawData[]> {
    const whitelistUpdates = await bridgeBank.queryFilter(bridgeBank.filters.LogWhiteListUpdate())
    return whitelistUpdates.map(x => {
        return {
            token: x.args._token,
            enabled: x.args._value,
            blockNumber: x.blockNumber,
            transactionIndex: x.transactionIndex
        }
    })
}

/**
 * Get whitelisted tokens and attach ERC20 data also.
 *
 * Use getWhitelistRawItems if you don't need the ERC20 data.
 *
 * @param bridgeBank
 * @param bridgeTokenFactory
 */
export async function getWhitelistItems(
    bridgeBank: BridgeBank,
    bridgeTokenFactory: BridgeToken__factory
): Promise<WhitelistItem[]> {
    const rawItems = await getWhitelistRawItems(bridgeBank)
    const xs = rawItems.map(async x => {
        const erc20Data = await getErc20Data(x.token, bridgeTokenFactory)
        return {
            erc20Data,
            logData: x
        }
    })
    return Promise.all(xs)
}

export async function getErc20Data(
    address: string,
    bridgeTokenFactory: BridgeToken__factory
): Promise<Erc20Data> {
    const bridgeToken = bridgeTokenFactory.attach(address)
    return {
        symbol: await bridgeToken.symbol(),
        decimals: await bridgeToken.decimals(),
        name: await bridgeToken.name(),
        address: bridgeToken.address
    }
}
