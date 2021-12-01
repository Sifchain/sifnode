package keeper

import (
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type KeeperI interface {
	InitGenesis(sdk.Context, types.GenesisState) []abci.ValidatorUpdate
	ExportGenesis(sdk.Context) *types.GenesisState

	SetMTP(ctx sdk.Context, mtp *types.MTP) error
	GetMTP(ctx sdk.Context, symbol string, mtpAddress string) (types.MTP, error)
	GetMTPIterator(ctx sdk.Context) sdk.Iterator
	GetMTPs(ctx sdk.Context) []*types.MTP
	GetMTPsForAsset(ctx sdk.Context, asset string) []*types.MTP
	GetAssetsForMTP(ctx sdk.Context, mtpAddress sdk.Address) []string
	DestroyMTP(ctx sdk.Context, symbol string, mtpAddress string) error
}
