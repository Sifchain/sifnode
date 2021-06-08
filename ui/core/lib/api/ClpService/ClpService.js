"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const entities_1 = require("../../entities");
const SifClient_1 = require("../utils/SifClient");
const toPool_1 = require("../utils/SifClient/toPool");
// TS not null type guard
function notNull(val) {
    return val !== null;
}
function createClpService({ sifApiUrl, nativeAsset, sifChainId, sifWsUrl, sifRpcUrl, sifUnsignedClient = new SifClient_1.SifUnSignedClient(sifApiUrl, sifWsUrl, sifRpcUrl), }) {
    const client = sifUnsignedClient;
    const instance = {
        getPools() {
            return __awaiter(this, void 0, void 0, function* () {
                try {
                    const rawPools = yield client.getPools();
                    return (rawPools
                        .map(toPool_1.toPool(nativeAsset))
                        // toPool can return a null pool for invalid pools lets filter them out
                        .filter(notNull));
                }
                catch (error) {
                    return [];
                }
            });
        },
        getPoolSymbolsByLiquidityProvider(address) {
            return __awaiter(this, void 0, void 0, function* () {
                // Unfortunately it is expensive for the backend to
                // filter pools so we need to annoyingly do this in two calls
                // First we get the metadata
                const poolMeta = yield client.getAssets(address);
                if (!poolMeta)
                    return [];
                return poolMeta.map(({ symbol }) => symbol);
            });
        },
        addLiquidity(params) {
            return __awaiter(this, void 0, void 0, function* () {
                return yield client.addLiquidity({
                    base_req: { chain_id: sifChainId, from: params.fromAddress },
                    external_asset: {
                        source_chain: params.externalAssetAmount.asset.network,
                        symbol: params.externalAssetAmount.asset.symbol,
                        ticker: params.externalAssetAmount.asset.symbol,
                    },
                    external_asset_amount: params.externalAssetAmount.toBigInt().toString(),
                    native_asset_amount: params.nativeAssetAmount.toBigInt().toString(),
                    signer: params.fromAddress,
                });
            });
        },
        createPool(params) {
            return __awaiter(this, void 0, void 0, function* () {
                return yield client.createPool({
                    base_req: { chain_id: sifChainId, from: params.fromAddress },
                    external_asset: {
                        source_chain: params.externalAssetAmount.asset.network,
                        symbol: params.externalAssetAmount.asset.symbol,
                        ticker: params.externalAssetAmount.asset.symbol,
                    },
                    external_asset_amount: params.externalAssetAmount.toBigInt().toString(),
                    native_asset_amount: params.nativeAssetAmount.toBigInt().toString(),
                    signer: params.fromAddress,
                });
            });
        },
        swap(params) {
            return __awaiter(this, void 0, void 0, function* () {
                return yield client.swap({
                    base_req: { chain_id: sifChainId, from: params.fromAddress },
                    received_asset: {
                        source_chain: params.receivedAsset.network,
                        symbol: params.receivedAsset.symbol,
                        ticker: params.receivedAsset.symbol,
                    },
                    sent_amount: params.sentAmount.toBigInt().toString(),
                    sent_asset: {
                        source_chain: params.sentAmount.asset.network,
                        symbol: params.sentAmount.asset.symbol,
                        ticker: params.sentAmount.asset.symbol,
                    },
                    min_receiving_amount: params.minimumReceived.toBigInt().toString(),
                    signer: params.fromAddress,
                });
            });
        },
        getLiquidityProvider(params) {
            return __awaiter(this, void 0, void 0, function* () {
                const response = yield client.getLiquidityProvider(params);
                let asset;
                const { LiquidityProvider: liquidityProvider, native_asset_balance, external_asset_balance, } = response.result;
                const { asset: { symbol }, liquidity_provider_units, liquidity_provider_address, } = liquidityProvider;
                try {
                    asset = entities_1.Asset(symbol);
                }
                catch (err) {
                    asset = entities_1.Asset({
                        name: symbol,
                        label: symbol,
                        symbol,
                        network: entities_1.Network.SIFCHAIN,
                        decimals: 18,
                    });
                }
                return entities_1.LiquidityProvider(asset, entities_1.Amount(liquidity_provider_units), liquidity_provider_address, entities_1.Amount(native_asset_balance), entities_1.Amount(external_asset_balance));
            });
        },
        removeLiquidity(params) {
            return __awaiter(this, void 0, void 0, function* () {
                return yield client.removeLiquidity({
                    asymmetry: params.asymmetry,
                    base_req: { chain_id: sifChainId, from: params.fromAddress },
                    external_asset: {
                        source_chain: params.asset.network,
                        symbol: params.asset.symbol,
                        ticker: params.asset.symbol,
                    },
                    signer: params.fromAddress,
                    w_basis_points: params.wBasisPoints,
                });
            });
        },
    };
    return instance;
}
exports.default = createClpService;
//# sourceMappingURL=ClpService.js.map