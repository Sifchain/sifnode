package types

import (
	adminkeeper "github.com/Sifchain/sifnode/x/admin/keeper"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
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

	CLPCalcSwap(ctx sdk.Context, sentAmount sdk.Uint, to clptypes.Asset, pool clptypes.Pool, marginEnabled bool) (sdk.Uint, error)

	GetPmtpRateParams(ctx sdk.Context) clptypes.PmtpRateParams

	GetRemovalQueue(ctx sdk.Context, symbol string) clptypes.RemovalQueue

	SingleExternalBalanceModuleAccountCheck(externalAsset string) sdk.Invariant
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
	GetMTPs(ctx sdk.Context, pagination *query.PageRequest) ([]*MTP, *query.PageResponse, error)
	GetMTPsForPool(ctx sdk.Context, asset string, pagination *query.PageRequest) ([]*MTP, *query.PageResponse, error)
	GetMTPsForAddress(ctx sdk.Context, mtpAddress sdk.Address, pagination *query.PageRequest) ([]*MTP, *query.PageResponse, error)
	DestroyMTP(ctx sdk.Context, mtpAddress string, id uint64) error

	IsWhitelisted(ctx sdk.Context, address string) bool
	WhitelistAddress(ctx sdk.Context, address string)
	DewhitelistAddress(ctx sdk.Context, address string)
	GetWhitelist(ctx sdk.Context, pagination *query.PageRequest) ([]string, *query.PageResponse, error)

	GetMaxLeverageParam(sdk.Context) sdk.Dec
	GetInterestRateMax(sdk.Context) sdk.Dec
	GetInterestRateMin(ctx sdk.Context) sdk.Dec
	GetInterestRateIncrease(ctx sdk.Context) sdk.Dec
	GetInterestRateDecrease(ctx sdk.Context) sdk.Dec
	GetHealthGainFactor(ctx sdk.Context) sdk.Dec
	GetEpochLength(ctx sdk.Context) int64
	GetPoolOpenThreshold(ctx sdk.Context) sdk.Dec
	GetRemovalQueueThreshold(ctx sdk.Context) sdk.Dec
	GetMaxOpenPositions(ctx sdk.Context) uint64
	GetEnabledPools(ctx sdk.Context) []string
	SetEnabledPools(ctx sdk.Context, pools []string)
	IsPoolEnabled(ctx sdk.Context, asset string) bool
	IsPoolClosed(ctx sdk.Context, asset string) bool
	IsWhitelistingEnabled(ctx sdk.Context) bool
	IsRowanCollateralEnabled(ctx sdk.Context) bool

	CLPSwap(ctx sdk.Context, sentAmount sdk.Uint, to string, pool clptypes.Pool) (sdk.Uint, error)
	Borrow(ctx sdk.Context, collateralAsset string, collateralAmount sdk.Uint, custodyAmount sdk.Uint, mtp *MTP, pool *clptypes.Pool, eta sdk.Dec) error
	TakeInCustody(ctx sdk.Context, mtp MTP, pool *clptypes.Pool) error
	TakeOutCustody(ctx sdk.Context, mtp MTP, pool *clptypes.Pool) error
	Repay(ctx sdk.Context, mtp *MTP, pool *clptypes.Pool, repayAmount sdk.Uint, takeFundPayment bool) error
	InterestRateComputation(ctx sdk.Context, pool clptypes.Pool) (sdk.Dec, error)
	CheckMinLiabilities(ctx sdk.Context, collateralAmount sdk.Uint, eta sdk.Dec, pool clptypes.Pool, custodyAsset string) error
	HandleInterestPayment(ctx sdk.Context, interestPayment sdk.Uint, mtp *MTP, pool *clptypes.Pool) sdk.Uint

	CalculatePoolHealth(pool *clptypes.Pool) sdk.Dec

	UpdatePoolHealth(ctx sdk.Context, pool *clptypes.Pool) error
	UpdateMTPHealth(ctx sdk.Context, mtp MTP, pool clptypes.Pool) (sdk.Dec, error)

	TrackSQBeginBlock(ctx sdk.Context, pool *clptypes.Pool)
	GetSQBeginBlock(ctx sdk.Context, pool *clptypes.Pool) uint64
	SetSQBeginBlock(ctx sdk.Context, pool *clptypes.Pool, height uint64)

	ForceCloseLong(ctx sdk.Context, mtp *MTP, pool *clptypes.Pool, isAdminClose bool, takeFundPayment bool) (sdk.Uint, error)

	EmitAdminClose(ctx sdk.Context, mtp *MTP, repayAmount sdk.Uint, closer string)
	EmitAdminCloseAll(ctx sdk.Context, takeMarginFund bool)

	GetSQFromQueue(ctx sdk.Context, pool clptypes.Pool) sdk.Dec
	GetSafetyFactor(ctx sdk.Context) sdk.Dec
}
