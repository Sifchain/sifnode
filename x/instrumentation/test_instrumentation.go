package instrumentation

import (
	"github.com/tendermint/tendermint/libs/log"
	"go.uber.org/zap"
)

const (
	// PeggyTestMarker is used in debugging logging statements for logs that are interesting for the Peggy test environment
	peggyTestMarker = "peggytest"
	KindMarker      = "kind"

	Startup       = "Startup"
	EthereumEvent = "EthereumEvent"
	CosmosEvent   = "CosmosEvent"
)

func PeggyCheckpoint(logger log.Logger, kind string, keysAndValues ...interface{}) {
	logger.Debug(peggyTestMarker, append([]interface{}{KindMarker, kind}, keysAndValues...)...)
}

func PeggyCheckpointZap(logger *zap.SugaredLogger, kind string, keysAndValues ...interface{}) {
	logger.Debugw(peggyTestMarker, append([]interface{}{KindMarker, kind}, keysAndValues...)...)
}
