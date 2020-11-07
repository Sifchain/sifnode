package clp

import (
	"github.com/Sifchain/sifnode/x/clp/keeper"
	types "github.com/Sifchain/sifnode/x/clp/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	NativeSymbol      = types.NativeSymbol
	NativeChain       = types.NativeChain
	NativeTicker      = types.NativeTicker
)

var (
	NewKeeper              = keeper.NewKeeper
	NewQuerier             = keeper.NewQuerier
	NewMsgSwap             = types.NewMsgSwap
	NewMsgAddLiquidity     = types.NewMsgAddLiquidity
	NewMsgRemoveLiquidity  = types.NewMsgRemoveLiquidity
	NewMsgCreatePool       = types.NewMsgCreatePool
	NewMsgDecommissionPool = types.NewMsgDecommissionPool
	NewLiquidityProvider   = types.NewLiquidityProvider
	NewAsset               = types.NewAsset
	NewPool                = types.NewPool
	RegisterCodec          = types.RegisterCodec
	NewGenesisState        = types.NewGenesisState
	DefaultGenesisState    = types.DefaultGenesisState
	NewParams              = types.NewParams
	ModuleCdc              = types.ModuleCdc
	CreateTestInputDefault = keeper.CreateTestInputDefault
	GenerateRandomPool     = keeper.GenerateRandomPool
	GenerateRandomLP       = keeper.GenerateRandomLP
	GenerateAddress        = keeper.GenerateAddress
	GenerateAddress2       = keeper.GenerateAddress2
	GetSettlementAsset     = types.GetSettlementAsset
)

type (
	Keeper              = keeper.Keeper
	MsgDecommissionPool = types.MsgDecommissionPool
	MsgCreatePool       = types.MsgCreatePool
	MsgAddLiquidity     = types.MsgAddLiquidity
	MsgRemoveLiquidity  = types.MsgRemoveLiquidity
	MsgSwap             = types.MsgSwap
	Pool                = types.Pool
	LiquidityProvider   = types.LiquidityProvider
	Asset               = types.Asset
	GenesisState        = types.GenesisState
)
