package faucet

import (
	"github.com/Sifchain/sifnode/x/faucet/internal/keeper"
	"github.com/Sifchain/sifnode/x/faucet/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	MAINNET = "mainnet"
	TESTNET = "testnet"
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper
)
