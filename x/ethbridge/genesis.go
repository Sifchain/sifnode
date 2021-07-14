package ethbridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/ethbridge/keeper"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
)

func DefaultGenesis() *types.GenesisState {
	return &types.GenesisState{}
}

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	// SetCrossChainFeeReceiverAccount
	if data.CrosschainFeeReceiveAccount != "" {
		receiveAccount, err := sdk.AccAddressFromBech32(data.CrosschainFeeReceiveAccount)
		if err != nil {
			panic(err)
		}
		keeper.SetCrossChainFeeReceiverAccount(ctx, receiveAccount)
	}

	// AddPeggyTokens
	if data.PeggyTokens != nil {
		for _, tokenStr := range data.PeggyTokens {
			keeper.AddPeggyToken(ctx, tokenStr)
		}
	}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	peggyTokens := keeper.GetPeggyToken(ctx)
	receiveAccount := keeper.GetCrossChainFeeReceiverAccount(ctx)

	return &types.GenesisState{
		PeggyTokens:                 peggyTokens.Tokens,
		CrosschainFeeReceiveAccount: receiveAccount.String(),
	}
}

func ValidateGenesis(data types.GenesisState) error {
	return nil
}
