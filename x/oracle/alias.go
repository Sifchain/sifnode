package oracle

// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/sifchain/peggy/x/oracle/keeper
// ALIASGEN: github.com/sifchain/peggy/x/oracle/types

import (
	"github.com/Sifchain/sifnode/x/oracle/keeper"
	"github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	DefaultConsensusNeeded = types.DefaultConsensusNeeded
	ModuleName             = types.ModuleName
	StoreKey               = types.StoreKey
	QuerierRoute           = types.QuerierRoute
	RouterKey              = types.RouterKey
	PendingStatusText      = types.PendingStatusText
	SuccessStatusText      = types.SuccessStatusText
	FailedStatusText       = types.FailedStatusText
)

var (
	// functions aliases
	NewKeeper                        = keeper.NewKeeper
	CreateTestAddrs                  = keeper.CreateTestAddrs
	CreateTestPubKeys                = keeper.CreateTestPubKeys
	CreateTestKeepers                = keeper.CreateTestKeepers
	NewClaim                         = types.NewClaim
	ErrProphecyNotFound              = types.ErrProphecyNotFound
	ErrMinimumConsensusNeededInvalid = types.ErrMinimumConsensusNeededInvalid
	ErrNoClaims                      = types.ErrNoClaims
	ErrInvalidIdentifier             = types.ErrInvalidIdentifier
	ErrProphecyFinalized             = types.ErrProphecyFinalized
	ErrDuplicateMessage              = types.ErrDuplicateMessage
	ErrInvalidClaim                  = types.ErrInvalidClaim
	ErrInvalidValidator              = types.ErrInvalidValidator
	ErrInternalDB                    = types.ErrInternalDB
	NewProphecy                      = types.NewProphecy
	NewStatus                        = types.NewStatus
	ModuleCdc                        = types.ModuleCdc
	GetGenesisStateFromAppState      = types.GetGenesisStateFromAppState
	NewNetworkDescriptor             = types.NewNetworkDescriptor

	// variable aliases
	StatusTextToString   = types.StatusTextToString
	StringToStatusText   = types.StringToStatusText
	MaxNetworkDescriptor = types.MaxNetworkDescriptor
)

type (
	Keeper             = keeper.Keeper
	Claim              = types.Claim
	Prophecy           = types.Prophecy
	DBProphecy         = types.DBProphecy
	Status             = types.Status
	StatusText         = types.StatusText
	GenesisState       = types.GenesisState
	NetworkDescriptor  = types.NetworkDescriptor
	ValidatorWhitelist = types.ValidatorWhitelist
)
