package ethbridge

// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/sifchain/sifnode/x/ethbridge/querier
// ALIASGEN: github.com/sifchain/sifnode/x/ethbridge/types

import (
	"github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

const (
	QueryEthProphecy = types.QueryEthProphecy
	ModuleName       = types.ModuleName
	StoreKey         = types.StoreKey
	QuerierRoute     = types.QuerierRoute
	RouterKey        = types.RouterKey
)

var (
	// functions aliases
	NewKeeper                   = keeper.NewKeeper
	NewQuerier                  = keeper.NewLegacyQuerier
	NewEthBridgeClaim           = types.NewEthBridgeClaim
	RegisterCodec               = types.RegisterLegacyAminoCodec
	ErrInvalidEthNonce          = types.ErrInvalidEthNonce
	ErrInvalidEthAddress        = types.ErrInvalidEthAddress
	ErrJSONMarshalling          = types.ErrJSONMarshalling
	NewEthereumAddress          = types.NewEthereumAddress
	NewMsgCreateEthBridgeClaim  = types.NewMsgCreateEthBridgeClaim
	NewQueryEthProphecyParams   = types.NewQueryEthProphecyRequest
	NewQueryEthProphecyResponse = types.NewQueryEthProphecyResponse

	CreateTestEthMsg                   = types.CreateTestEthMsg
	CreateTestEthClaim                 = types.CreateTestEthClaim
	CreateTestQueryEthProphecyResponse = types.CreateTestQueryEthProphecyResponse
	CrossChainFeeReceiverAccountPrefix = types.CrossChainFeeReceiverAccountPrefix
)

type (
	Keeper                                = keeper.Keeper
	EthBridgeClaim               = types.EthBridgeClaim //nolint:revive
	OracleClaimContent           = types.OracleClaimContent
	EthereumAddress                       = types.EthereumAddress
	MsgCreateEthBridgeClaim               = types.MsgCreateEthBridgeClaim
	MsgBurn                               = types.MsgBurn
	MsgLock                               = types.MsgLock
	MsgUpdateWhiteListValidator           = types.MsgUpdateWhiteListValidator
	MsgUpdateCrossChainFeeReceiverAccount = types.MsgUpdateCrossChainFeeReceiverAccount
	MsgRescueCrossChainFee                = types.MsgRescueCrossChainFee
	MsgSetFeeInfo                         = types.MsgSetFeeInfo
	MsgSignProphecy                       = types.MsgSignProphecy
	QueryEthProphecyParams                = types.QueryEthProphecyRequest
	QueryEthProphecyResponse              = types.QueryEthProphecyResponse
	MsgUpdateConsensusNeeded              = types.MsgUpdateConsensusNeeded
)
