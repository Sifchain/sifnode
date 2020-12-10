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
	MaxWbasis         = types.MaxWbasis
)

var (
	NewKeeper                   = keeper.NewKeeper
	NewQuerier                  = keeper.NewQuerier
	NewMsgSwap                  = types.NewMsgSwap
	NewMsgAddLiquidity          = types.NewMsgAddLiquidity
	NewMsgRemoveLiquidity       = types.NewMsgRemoveLiquidity
	NewMsgCreatePool            = types.NewMsgCreatePool
	NewMsgDecommissionPool      = types.NewMsgDecommissionPool
	NewAsset                    = types.NewAsset
	RegisterCodec               = types.RegisterCodec
	DefaultGenesisState         = types.DefaultGenesisState
	ModuleCdc                   = types.ModuleCdc
	GetSettlementAsset          = types.GetSettlementAsset
	GetGenesisStateFromAppState = types.GetGenesisStateFromAppState
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
