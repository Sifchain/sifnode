package v42

import (
	v039ethbridge "github.com/Sifchain/sifnode/x/ethbridge/legacy/v39"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

func Migrate(state v039ethbridge.GenesisState) *types.GenesisState {
	return &types.GenesisState{
		CrosschainFeeReceiveAccount: state.CethReceiverAccount.String(),
		PeggyTokens:                 state.PeggyTokens,
	}
}
