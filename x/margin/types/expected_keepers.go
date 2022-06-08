package types

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/abci/types"
)

type BankKeeper interface {
	//SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
}

type CLPKeeper interface {
	GetPools(ctx sdk.Context) []*clptypes.Pool
	GetPool(ctx sdk.Context, symbol string) (clptypes.Pool, error)
	SetPool(ctx sdk.Context, pool *clptypes.Pool) error
	GetNormalizationFactorFromAsset(sdk.Context, clptypes.Asset) (sdk.Dec, bool, error)

	ValidateZero(inputs []sdk.Uint) bool
	ReducePrecision(dec sdk.Dec, po int64) sdk.Dec
	IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec
	GetMinLen(inputs []sdk.Uint) int64

	GetPmtpRateParams(ctx sdk.Context) clptypes.PmtpRateParams
}

type Keeper interface {
	InitGenesis(sdk.Context, GenesisState) []types.ValidatorUpdate
	ExportGenesis(sdk.Context) *GenesisState
	BeginBlocker(sdk.Context)

	ClpKeeper() CLPKeeper
	BankKeeper() BankKeeper

	SetMTP(ctx sdk.Context, mtp *MTP) error
	GetMTP(ctx sdk.Context, address string, id uint64) (MTP, error)
	GetMTPIterator(ctx sdk.Context) sdk.Iterator
	GetMTPs(ctx sdk.Context) []*MTP
	GetMTPsForCollateralAsset(ctx sdk.Context, asset string) []*MTP
	GetMTPsForCustodyAsset(ctx sdk.Context, asset string) []*MTP
	GetAssetsForMTP(ctx sdk.Context, mtpAddress sdk.Address) []string
	GetMTPsForAddress(ctx sdk.Context, mtpAddress sdk.Address) []*MTP
	DestroyMTP(sdk.Context, string, uint64) error

	GetLeverageParam(sdk.Context) sdk.Uint
	GetInterestRateMax(sdk.Context) sdk.Dec
	GetInterestRateMin(ctx sdk.Context) sdk.Dec
	GetInterestRateIncrease(ctx sdk.Context) sdk.Dec
	GetInterestRateDecrease(ctx sdk.Context) sdk.Dec
	GetHealthGainFactor(ctx sdk.Context) sdk.Dec
	GetEpochLength(ctx sdk.Context) int64
	GetForceCloseThreshold(ctx sdk.Context) sdk.Dec
	GetEnabledPools(ctx sdk.Context) []string
	SetEnabledPools(ctx sdk.Context, pools []string)
	IsPoolEnabled(ctx sdk.Context, asset string) bool

	CustodySwap(ctx sdk.Context, pool clptypes.Pool, to string, sentAmount sdk.Uint) (sdk.Uint, error)
	Borrow(ctx sdk.Context, collateralAsset string, collateralAmount sdk.Uint, borrowAmount sdk.Uint, mtp *MTP, pool *clptypes.Pool, leverage sdk.Uint) error
	TakeInCustody(ctx sdk.Context, mtp MTP, pool *clptypes.Pool) error
	TakeOutCustody(ctx sdk.Context, mtp MTP, pool *clptypes.Pool) error
	Repay(ctx sdk.Context, mtp *MTP, pool clptypes.Pool, repayAmount sdk.Uint) error
	InterestRateComputation(ctx sdk.Context, pool clptypes.Pool) (sdk.Dec, error)

	UpdateMTPInterestLiabilities(ctx sdk.Context, mtp *MTP, interestRate sdk.Dec) error
	UpdatePoolHealth(ctx sdk.Context, pool *clptypes.Pool) error
	UpdateMTPHealth(ctx sdk.Context, mtp MTP, pool clptypes.Pool) (sdk.Dec, error)

	ForceCloseLong(ctx sdk.Context, msg *MsgForceClose) (*MTP, error)
}
