package faucet

import (
	"github.com/Sifchain/sifnode/x/faucet/keeper"
	types "github.com/Sifchain/sifnode/x/faucet/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	TESTNET           = types.TESTNET
	MAINNET           = types.MAINNET
)

var (
	NewKeeper                   = keeper.NewKeeper
	NewQuerier                  = keeper.NewQuerier
	GetFaucetModuleAddress      = types.GetFaucetModuleAddress
	NewMsgRequestCoins          = types.NewMsgRequestCoins
	NewMsgAddCoins              = types.NewMsgAddCoins
	NewRemoveFacuetFeeDecorator = types.NewRemoveFacuetFeeDecorator
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)
