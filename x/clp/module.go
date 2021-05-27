package clp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/clp/client/cli"
	"github.com/Sifchain/sifnode/x/clp/client/rest"
	"github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/types"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the clp module.
type AppModuleBasic struct{}

// Name returns the clp module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the clp module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (b AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the clp
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the clp module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the clp module.
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the clp module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic("Failed to register GRPC gateway routes.")
	}
}

// GetTxCmd returns the root tx command for the clp module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command for the clp module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

//____________________________________________________________________________

// AppModule implements an application module for the clp module.
type AppModule struct {
	AppModuleBasic

	keeper              keeper.Keeper
	bankKeeper          types.BankKeeper
	tokenMap            map[string]string
	migrationStartBlock int64
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, bankKeeper types.BankKeeper, tokenMap map[string]string, migrationStartBlock int64) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		keeper:              k,
		bankKeeper:          bankKeeper,
		tokenMap:            tokenMap,
		migrationStartBlock: migrationStartBlock,
	}
}

// Name returns the clp module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the clp module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

// QuerierRoute returns the clp module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// LegacyQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	querier := keeper.Querier{Keeper: am.keeper}
	types.RegisterQueryServer(cfg.QueryServer(), querier)
}

// InitGenesis performs genesis initialization for the clp module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the clp
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)

	return cdc.MustMarshalJSON(&gs)
}

// BeginBlock used to do token migration
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	if req.Header.Height == am.migrationStartBlock {
		am.migrateBalance(ctx, am.tokenMap)
		am.migrateLiquidity(ctx, am.tokenMap)
		am.migratePool(ctx, am.tokenMap)
	}
}

func getAll(addresses *[]sdk.AccAddress, coins *[]sdk.Coin) func(address sdk.AccAddress, coin sdk.Coin) bool {
	return func(address sdk.AccAddress, coin sdk.Coin) bool {
		*addresses = append(*addresses, address)
		*coins = append(*coins, coin)
		return true
	}
}

func (am AppModule) migrateBalance(ctx sdk.Context, tokenMap map[string]string) {
	addresses := []sdk.AccAddress{}
	coins := []sdk.Coin{}

	am.bankKeeper.IterateAllBalances(ctx, getAll(&addresses, &coins))

	for index, address := range addresses {

		coin := coins[index]
		amount := coin.Amount
		// clear the balance for old denom
		coin.Amount = sdk.NewInt(0)
		err := am.bankKeeper.SetBalance(ctx, address, coin)
		if err != nil {
			panic("failed to set balance during token migration")
		}

		// set the balance for new denom
		if value, ok := tokenMap[coin.Denom]; ok {
			coin = sdk.NewCoin(value, amount)
			err = am.bankKeeper.SetBalance(ctx, address, coin)
			if err != nil {
				panic("failed to set balance during token migration")
			}
		} else {
			panic(fmt.Sprintf("new denom for %s not found\n", coin.Denom))
		}
	}

}

func (am AppModule) migratePool(ctx sdk.Context, tokenMap map[string]string) {
	pools := am.keeper.GetPools(ctx)
	for _, value := range pools {
		token := value.ExternalAsset.Symbol
		if newDenom, ok := tokenMap[token]; ok {
			err := am.keeper.DestroyPool(ctx, token)
			if err != nil {
				panic("failed to destroy pool during token migration")
			}
			value.ExternalAsset.Symbol = newDenom
			err = am.keeper.SetPool(ctx, value)
			if err != nil {
				panic("failed to set pool during token migration")
			}

		} else {
			panic(fmt.Sprintf("new denom for %s not found\n", token))
		}
	}
}

func (am AppModule) migrateLiquidity(ctx sdk.Context, tokenMap map[string]string) {
	liquidity := am.keeper.GetLiquidityProviders(ctx)

	for _, value := range liquidity {
		token := value.Asset.Symbol
		if newDenom, ok := tokenMap[token]; ok {

			am.keeper.DestroyLiquidityProvider(ctx, token, value.LiquidityProviderAddress)

			value.Asset.Symbol = newDenom
			am.keeper.SetLiquidityProvider(ctx, value)

		} else {
			panic(fmt.Sprintf("new denom for %s not found\n", token))
		}
	}
}

// EndBlock returns the end blocker for the clp module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
