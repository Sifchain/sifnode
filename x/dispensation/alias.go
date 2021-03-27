package dispensation

import (
	"github.com/Sifchain/sifnode/x/dispensation/keeper"
	types "github.com/Sifchain/sifnode/x/dispensation/types"
)

const (
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState
	ModuleCdc           = types.ModuleCdc
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	MsgAirdrop   = types.MsgAirdrop
)
