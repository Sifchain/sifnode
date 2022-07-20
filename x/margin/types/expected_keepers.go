//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

import (
	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/abci/types"
)

type BankKeeper interface {
	//SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

type CLPKeeper interface {
	GetPools(ctx sdk.Context) []*clptypes.Pool
	GetPool(ctx sdk.Context, symbol string) (clptypes.Pool, error)
	SetPool(ctx sdk.Context, pool *clptypes.Pool) error
	GetNormalizationFactorFromAsset(ctx sdk.Context, asset clptypes.Asset) (sdk.Dec, bool, error)

	ValidateZero(inputs []sdk.Uint) bool
	ReducePrecision(dec sdk.Dec, po int64) sdk.Dec
	IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec
	GetMinLen(inputs []sdk.Uint) int64

	GetPmtpRateParams(ctx sdk.Context) clptypes.PmtpRateParams

	GetRemovalQueue(ctx sdk.Context, symbol string) clptypes.RemovalQueue
}

type Keeper interface {
	InitGenesis(ctx sdk.Context, data GenesisState) []types.ValidatorUpdate
	ExportGenesis(sdk.Context) *GenesisState
	BeginBlocker(sdk.Context)

	ClpKeeper() CLPKeeper
	BankKeeper() BankKeeper
	AdminKeeper() adminkeeper.Keeper

	GetParams(sdk.Context) Params
	SetParams(sdk.Context, *Params)

	GetMTPCount(ctx sdk.Context) uint64
	GetOpenMTPCount(ctx sdk.Context) uint64
	SetMTP(ctx sdk.Context, mtp *MTP) error
	GetMTP(ctx sdk.Context, mtpAddress string, id uint64) (MTP, error)
	GetMTPIterator(ctx sdk.Context) sdk.Iterator
	GetMTPs(ctx sdk.Context) []*MTP
	GetMTPsForPool(ctx sdk.Context, asset string) []*MTP
	GetAssetsForMTP(ctx sdk.Context, mtpAddress sdk.Address) []string
	GetMTPsForAddress(ctx sdk.Context, mtpAddress sdk.Address) []*MTP
	DestroyMTP(ctx sdk.Context, mtpAddress string, id uint64) error

	GetMaxLeverageParam(sdk.Context) sdk.Dec
	GetInterestRateMax(sdk.Context) sdk.Dec
	GetInterestRateMin(ctx sdk.Context) sdk.Dec
	GetInterestRateIncrease(ctx sdk.Context) sdk.Dec
	GetInterestRateDecrease(ctx sdk.Context) sdk.Dec
	GetHealthGainFactor(ctx sdk.Context) sdk.Dec
	GetEpochLength(ctx sdk.Context) int64
	GetForceCloseThreshold(ctx sdk.Context) sdk.Dec
	GetRemovalQueueThreshold(ctx sdk.Context) sdk.Dec
	GetMaxOpenPositions(ctx sdk.Context) uint64
	GetEnabledPools(ctx sdk.Context) []string
	SetEnabledPools(ctx sdk.Context, pools []string)
	IsPoolEnabled(ctx sdk.Context, asset string) bool

	CustodySwap(ctx sdk.Context, pool clptypes.Pool, to string, sentAmount sdk.Uint) (sdk.Uint, error)
	Borrow(ctx sdk.Context, collateralAsset string, collateralAmount sdk.Uint, custodyAmount sdk.Uint, mtp *MTP, pool *clptypes.Pool, eta sdk.Dec) error
	TakeInCustody(ctx sdk.Context, mtp MTP, pool *clptypes.Pool) error
	TakeOutCustody(ctx sdk.Context, mtp MTP, pool *clptypes.Pool) error
	Repay(ctx sdk.Context, mtp *MTP, pool clptypes.Pool, repayAmount sdk.Uint) error
	InterestRateComputation(ctx sdk.Context, pool clptypes.Pool) (sdk.Dec, error)

	CalculatePoolHealth(pool *clptypes.Pool) sdk.Dec

	UpdateMTPInterestLiabilities(ctx sdk.Context, mtp *MTP, interestRate sdk.Dec) error
	UpdatePoolHealth(ctx sdk.Context, pool *clptypes.Pool) error
	UpdateMTPHealth(ctx sdk.Context, mtp MTP, pool clptypes.Pool) (sdk.Dec, error)

	ForceCloseLong(ctx sdk.Context, msg *MsgForceClose) (*MTP, error)

	EmitForceClose(ctx sdk.Context, mtp *MTP, closer string)

	GetSQ(ctx sdk.Context, pool clptypes.Pool) sdk.Dec
}
