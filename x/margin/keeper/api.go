package keeper

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	"github.com/Sifchain/sifnode/x/margin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type KeeperI interface {
	InitGenesis(sdk.Context, types.GenesisState) []abci.ValidatorUpdate
	ExportGenesis(sdk.Context) *types.GenesisState

	ClpKeeper() types.CLPKeeper
	BankKeeper() types.BankKeeper

	SetMTP(ctx sdk.Context, mtp *types.MTP) error
	GetMTP(ctx sdk.Context, symbol string, mtpAddress string) (types.MTP, error)
	GetMTPIterator(ctx sdk.Context) sdk.Iterator
	GetMTPs(ctx sdk.Context) []*types.MTP
	GetMTPsForAsset(ctx sdk.Context, asset string) []*types.MTP
	GetAssetsForMTP(ctx sdk.Context, mtpAddress sdk.Address) []string
	DestroyMTP(ctx sdk.Context, symbol string, mtpAddress string) error

	GetLeverageParam(sdk.Context) sdk.Uint

	CustodySwap(ctx sdk.Context, pool clptypes.Pool, to string, sentAmount sdk.Uint) (sdk.Uint, error)
	Borrow(ctx sdk.Context, collateralAsset string, collateralAmount sdk.Uint, borrowAmount sdk.Uint, mtp types.MTP, pool clptypes.Pool, leverage sdk.Uint) error
	TakeInCustody(ctx sdk.Context, mtp types.MTP, pool clptypes.Pool) error

	UpdatePoolHealth(ctx sdk.Context, pool clptypes.Pool) error
	UpdateMTPHealth(ctx sdk.Context, mtp types.MTP, pool clptypes.Pool) (sdk.Dec, error)
}
