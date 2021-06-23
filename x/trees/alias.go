package trees

import (
	"github.com/Sifchain/sifnode/x/trees/keeper"
	types "github.com/Sifchain/sifnode/x/trees/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute

	MAINNET = "mainnet"
	TESTNET = "testnet"
)

var (
	NewKeeper              = keeper.NewKeeper
	NewQuerier             = keeper.NewQuerier
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)
