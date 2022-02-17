package instrumentation

import (
	"github.com/tendermint/tendermint/libs/log"
	"go.uber.org/zap"
)

const (
	// PeggyTestMarker is used in debugging logging statements for logs that are interesting for the Peggy test environment
	peggyTestMarker = "peggytest"
	kindMarker      = "kind"

	Startup       = "Startup"
	EthereumEvent = "EthereumEvent"
	CosmosEvent   = "CosmosEvent"
	// TODO: rename, there isnt ethereumprophecyclaim, is there?
	EthereumProphecyClaim          = "EthereumProphecyClaim"
	SetGlobalSequenceToBlockNumber = "SetGlobalSequenceToBlockNumber"
	SendCoinsFromAccountToModule   = "SendCoinsFromAccountToModule"
	BurnCoins                      = "BurnCoins"
	SignProphecy                   = "SignProphecy"
	ProcessSuccessfulClaim         = "ProcessSuccessfulClaim"
	CoinsSent                      = "coinsSent"
	Burn                           = "Burn"
	CreateEthBridgeClaim           = "CreateEthBridgeClaim"
	Lock                           = "Lock"
	GetCrossChainFeeConfig         = "GetCrossChainFeeConfig"
	AppendValidatorToProphecy      = "AppendValidatorToProphecy"
	ProcessCompletion              = "processCompletion"
	ProphecyStatus                 = "ProphecyStatus"
	AppendSignature                = "AppendSignature"
	SetGlobalNonceProphecyID       = "SetGlobalNonceProphecyID"
	SetProphecy                    = "SetProphecy"
	SetProphecyInfo                = "SetProphecyInfo"
	SetWitnessLockBurnNonce        = "SetWitnessLockBurnNonce"
	SetFirstLockDoublePeg          = "SetFirstLockDoublePeg"
	AddTokenMetadata               = "AddTokenMetadata"
	GetTokenMetadata               = "GetTokenMetadata"
	PublishCosmosLockMessage       = "PublishCosmosLockMessage"
	PublishCosmosBurnMessage       = "PublishCosmosBurnMessage"
	ReceiveCosmosLockMessage       = "ReceiveCosmosLockMessage"
	ReceiveCosmosBurnMessage       = "ReceiveCosmosBurnMessage"
	WitnessSignProphecy            = "WitnessSignProphecy"
	// TODO: Ambiguous naming
	ProphecyClaimSubmitted = "ProphecyClaimSubmitted"
)

func PeggyCheckpoint(logger log.Logger, kind string, keysAndValues ...interface{}) {
	logger.Debug(peggyTestMarker, append([]interface{}{kindMarker, kind}, keysAndValues...)...)
}

func PeggyCheckpointZap(logger *zap.SugaredLogger, kind string, keysAndValues ...interface{}) {
	logger.Debugw(peggyTestMarker, append([]interface{}{kindMarker, kind}, keysAndValues...)...)
}
