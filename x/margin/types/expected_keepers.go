package types

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/abci/types"
)

type BankKeeper interface {
	//SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	//SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type CLPKeeper interface {
	GetPool(ctx sdk.Context, symbol string) (clptypes.Pool, error)
	SetPool(ctx sdk.Context, pool *clptypes.Pool) error
	GetNormalizationFactorForAsset(sdk.Context, string) (sdk.Dec, bool, error)

	ValidateZero(inputs []sdk.Uint) bool
	ReducePrecision(dec sdk.Dec, po int64) sdk.Dec
	IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec
	GetMinLen(inputs []sdk.Uint) int64
}

type Keeper interface {
	InitGenesis(sdk.Context, GenesisState) []types.ValidatorUpdate
	ExportGenesis(sdk.Context) *GenesisState

	ClpKeeper() CLPKeeper
	BankKeeper() BankKeeper

	SetMTP(ctx sdk.Context, mtp *MTP) error
	GetMTP(ctx sdk.Context, symbol string, mtpAddress string) (MTP, error)
	GetMTPIterator(ctx sdk.Context) sdk.Iterator
	GetMTPs(ctx sdk.Context) []*MTP
	GetMTPsForAsset(ctx sdk.Context, asset string) []*MTP
	GetAssetsForMTP(ctx sdk.Context, mtpAddress sdk.Address) []string
	DestroyMTP(ctx sdk.Context, symbol string, mtpAddress string) error

	GetLeverageParam(sdk.Context) sdk.Uint

	CustodySwap(ctx sdk.Context, pool clptypes.Pool, to string, sentAmount sdk.Uint) (sdk.Uint, error)
	Borrow(ctx sdk.Context, collateralAsset string, collateralAmount sdk.Uint, borrowAmount sdk.Uint, mtp MTP, pool clptypes.Pool, leverage sdk.Uint) error
	TakeInCustody(ctx sdk.Context, mtp MTP, pool clptypes.Pool) error

	UpdatePoolHealth(ctx sdk.Context, pool clptypes.Pool) error
	UpdateMTPHealth(ctx sdk.Context, mtp MTP, pool clptypes.Pool) (sdk.Dec, error)
}
