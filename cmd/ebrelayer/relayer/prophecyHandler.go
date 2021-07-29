package relayer

import (
	cosmosbridge "github.com/Sifchain/sifnode/cmd/ebrelayer/contract/generated/bindings/cosmosbridge"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
)

var currentGlobalNonce = uint64(0)
var prophecyInfoBuffer = map[uint64]types.ProphecyInfo{}

// Parses event data from the msg, event, builds a new ProphecyClaim, and relays it to Ethereum
func (sub CosmosSub) handleProphecyCompleted(
	prophecyInfo types.ProphecyInfo,
	claimType types.Event,
) {
	sub.SugaredLogger.Infow(
		"get the prophecy completed message.",
		"cosmosMsg", prophecyInfo,
	)

	// will discard it if global nonce is less than current global nonce
	if currentGlobalNonce > prophecyInfo.GlobalNonce {
		sub.SugaredLogger.Errorw("global nonce is invalid.",
			"current global nonce", currentGlobalNonce,
			"global nonce in message", prophecyInfo.GlobalNonce)
		return
	}

	// buffer the prophecy
	if prophecyInfo.GlobalNonce > currentGlobalNonce+1 {
		prophecyInfoBuffer[prophecyInfo.GlobalNonce] = prophecyInfo
		return
	}

	client, auth, target, err := tryInitRelayConfig(sub)
	if err != nil {
		sub.SugaredLogger.Errorw("failed in init relay config.",
			errorMessageKey, err.Error())
		return
	}

	// Initialize CosmosBridge instance
	cosmosBridgeInstance, err := cosmosbridge.NewCosmosBridge(target, client)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get cosmosBridge instance.",
			errorMessageKey, err.Error())
		return
	}

	maxRetries := 5
	i := 0

	for i < maxRetries {
		err = txs.RelayProphecyCompletedToEthereum(
			prophecyInfo,
			sub.SugaredLogger,
			client,
			auth,
			cosmosBridgeInstance,
		)

		if err != nil {
			sub.SugaredLogger.Errorw(
				"failed to send new prophecy completed to ethereum",
				errorMessageKey, err.Error(),
			)
		} else {
			break
		}
		i++
	}

	if i == maxRetries {
		sub.SugaredLogger.Errorw(
			"failed to broadcast transaction after 5 attempts",
			errorMessageKey, err.Error(),
		)
	}

	// update currentGlobalNonce
	currentGlobalNonce++
}
