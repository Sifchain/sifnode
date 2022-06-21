package instrumentation

import (
	"os"
	"strings"

	"github.com/tendermint/tendermint/libs/log"
	"go.uber.org/zap"
)

const (
	// PeggyTestMarker is used in debugging logging statements for logs that are interesting for the Peggy test environment
	peggyTestMarker = "peggytest"
	kindMarker      = "kind"

	Startup                        = "Startup"
	EthereumEvent                  = "EthereumEvent"
	CosmosEvent                    = "CosmosEvent"
	EthereumBridgeClaim            = "EthereumBridgeClaim"
	SetGlobalSequenceToBlockNumber = "SetGlobalSequenceToBlockNumber"
	SendCoinsFromAccountToModule   = "SendCoinsFromAccountToModule"
	BurnCoins                      = "BurnCoins"
	SignProphecy                   = "SignProphecy"
	ProcessSignProphecy            = "ProcessSignProphecy"
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
	SetFirstDoublePeg              = "SetFirstDoublePeg"
	AddTokenMetadata               = "AddTokenMetadata"
	GetTokenMetadata               = "GetTokenMetadata"
	PublishCosmosBurnMessage       = "PublishCosmosBurnMessage"
	ReceiveCosmosBurnMessage       = "ReceiveCosmosBurnMessage"
	WitnessSignProphecy            = "WitnessSignProphecy"
	ProphecyClaimSubmitted         = "ProphecyClaimSubmitted"
)

func PeggyCheckpoint(logger log.Logger, kind string, keysAndValues ...interface{}) {
	logger.Debug(peggyTestMarker, append([]interface{}{kindMarker, kind, "cmdline", strings.Join(os.Args, " ")}, keysAndValues...)...)
}

func PeggyCheckpointZap(logger *zap.SugaredLogger, kind string, keysAndValues ...interface{}) {
	logger.Debugw(peggyTestMarker, append([]interface{}{kindMarker, kind, "cmdline", strings.Join(os.Args, " ")}, keysAndValues...)...)
}
